// bench/main.go
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/examples"
	"github.com/tuneinsight/lattigo/v6/schemes/bgv"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

var (
	degrees = []int{4096, 8192, 16384}
	ops     = []string{
		"Генерация релинеаризации",
		"Генерация Galois ключей",
		"Шифрование",
		"Дешифрование",
		"Сложение",
		"Умножение",
		"Релионеаризация",
		"Вращение (rotate rows 1 шаг)",
		"Сериализация шифротекста",
		"Сжатая сериализация (ZLIB)",
		"Сжатая сериализация (Zstandard)",
	}
)

func main() {
	outFile := "results.csv"
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("cannot create output csv: %v", err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	// Header in your template style (Russian)
	header := []string{"Операция"}
	for _, sch := range []string{"BFV", "CKKS", "BGV"} {
		for _, d := range degrees {
			header = append(header, fmt.Sprintf("%s %d", sch, d))
		}
	}
	if err := w.Write(header); err != nil {
		log.Fatalf("csv write header: %v", err)
	}

	// For each operation we'll collect times per scheme+degree
	// Build ordering: BFV-4096 BFV-8192 BFV-16384 CKKS-4096 ...
	resultGrid := make(map[string]map[string]float64) // op -> (schemeDeg -> ms)
	for _, op := range ops {
		resultGrid[op] = map[string]float64{}
	}

	// ---- BGV ----
	// Примечание: BFV в Lattigo v6 не имеет отдельной схемы, используется BGV
	for _, d := range degrees {
		key := fmt.Sprintf("BFV %d", d)
		fmt.Println("Running BGV (BFV)", d)
		var pl bgv.ParametersLiteral
		switch d {
		case 4096:
			pl = examples.BGVParamsN12QP109
		case 8192:
			pl = examples.BGVParamsN13QP218
		case 16384:
			pl = examples.BGVParamsN14QP438
		default:
			log.Fatalf("unsupported degree: %d", d)
		}
		params, err := bgv.NewParametersFromLiteral(pl)
		if err != nil {
			log.Fatalf("failed to create BGV params: %v", err)
		}
		benchBGV(params, key, resultGrid)
	}

	// ---- CKKS ----
	for _, d := range degrees {
		key := fmt.Sprintf("CKKS %d", d)
		fmt.Println("Running CKKS", d)
		var pl ckks.ParametersLiteral
		switch d {
		case 4096:
			pl = examples.CKKSComplexParamsN12QP109
		case 8192:
			pl = examples.CKKSComplexParamsN13QP218
		case 16384:
			pl = examples.CKKSComplexParamsN14QP438
		default:
			log.Fatalf("unsupported degree: %d", d)
		}
		params, err := ckks.NewParametersFromLiteral(pl)
		if err != nil {
			log.Fatalf("failed to create CKKS params: %v", err)
		}
		benchCKKS(params, key, resultGrid)
	}

	// ---- BGV ----
	for _, d := range degrees {
		key := fmt.Sprintf("BGV %d", d)
		fmt.Println("Running BGV", d)
		var pl bgv.ParametersLiteral
		switch d {
		case 4096:
			pl = examples.BGVParamsN12QP109
		case 8192:
			pl = examples.BGVParamsN13QP218
		case 16384:
			pl = examples.BGVParamsN14QP438
		default:
			log.Fatalf("unsupported degree: %d", d)
		}
		params, err := bgv.NewParametersFromLiteral(pl)
		if err != nil {
			log.Fatalf("failed to create BGV params: %v", err)
		}
		benchBGV(params, key, resultGrid)
	}

	// Write CSV rows in the order ops[] and columns as header
	for _, op := range ops {
		row := []string{op}
		for _, sch := range []string{"BFV", "CKKS", "BGV"} {
			for _, d := range degrees {
				k := fmt.Sprintf("%s %d", sch, d)
				val := resultGrid[op][k]
				row = append(row, fmt.Sprintf("%.3f", val))
			}
		}
		w.Write(row)
	}
	fmt.Println("Done. CSV:", filepath.Join(".", outFile))
}

// benchBGV выполняет бенчмарки операций схемы BGV
func benchBGV(params bgv.Parameters, key string, grid map[string]map[string]float64) {
	N := 100 // Количество итераций для усреднения

	// Создаем ключи один раз
	kgen := rlwe.NewKeyGenerator(params)
	sk := kgen.GenSecretKeyNew()

	// ---- Генерация релинеаризации ----
	durRelinKeyGen := time.Duration(0)
	var rlk *rlwe.RelinearizationKey
	for i := 0; i < N; i++ {
		start := time.Now()
		rlk = kgen.GenRelinearizationKeyNew(sk)
		durRelinKeyGen += time.Since(start)
	}
	grid["Генерация релинеаризации"][key] = float64(durRelinKeyGen.Microseconds()) / float64(N)

	// ---- Генерация Galois ключей ----
	durGaloisKeyGen := time.Duration(0)
	galEls := params.GaloisElementForRowRotation()
	for i := 0; i < N; i++ {
		start := time.Now()
		kgen.GenGaloisKeyNew(galEls, sk)
		durGaloisKeyGen += time.Since(start)
	}
	grid["Генерация Galois ключей"][key] = float64(durGaloisKeyGen.Microseconds()) / float64(N)

	// Создаем Galois ключ для вращений
	gk := kgen.GenGaloisKeyNew(galEls, sk)
	evk := rlwe.NewMemEvaluationKeySet(rlk, gk)

	// Подготовка объектов
	encoder := bgv.NewEncoder(params)
	encryptor := rlwe.NewEncryptor(params, sk)
	decryptor := rlwe.NewDecryptor(params, sk)
	evaluator := bgv.NewEvaluator(params, evk)

	// Создание тестовых данных
	values1 := make([]uint64, params.MaxSlots())
	values2 := make([]uint64, params.MaxSlots())
	T := params.PlaintextModulus()
	r := rand.New(rand.NewSource(0))
	for i := range values1 {
		values1[i] = r.Uint64() % T
		values2[i] = r.Uint64() % T
	}

	pt1 := bgv.NewPlaintext(params, params.MaxLevel())
	pt2 := bgv.NewPlaintext(params, params.MaxLevel())
	encoder.Encode(values1, pt1)
	encoder.Encode(values2, pt2)

	// ---- Шифрование ----
	ct1, _ := encryptor.EncryptNew(pt1)
	ct2, _ := encryptor.EncryptNew(pt2)

	durEncrypt := time.Duration(0)
	for i := 0; i < N; i++ {
		start := time.Now()
		encryptor.EncryptNew(pt1)
		durEncrypt += time.Since(start)
	}
	grid["Шифрование"][key] = float64(durEncrypt.Microseconds()) / float64(N)

	// ---- Дешифрование ----
	durDecrypt := time.Duration(0)
	for i := 0; i < N; i++ {
		start := time.Now()
		decryptor.DecryptNew(ct1)
		durDecrypt += time.Since(start)
	}
	grid["Дешифрование"][key] = float64(durDecrypt.Microseconds()) / float64(N)

	// ---- Сложение ----
	durAdd := time.Duration(0)
	ctAdd := bgv.NewCiphertext(params, 1, params.MaxLevel())
	for i := 0; i < N; i++ {
		start := time.Now()
		evaluator.Add(ct1, ct2, ctAdd)
		durAdd += time.Since(start)
	}
	grid["Сложение"][key] = float64(durAdd.Microseconds()) / float64(N)

	// ---- Умножение ----
	durMul := time.Duration(0)
	ctMul := bgv.NewCiphertext(params, 2, params.MaxLevel())
	for i := 0; i < N; i++ {
		start := time.Now()
		evaluator.Mul(ct1, ct2, ctMul)
		durMul += time.Since(start)
	}
	grid["Умножение"][key] = float64(durMul.Microseconds()) / float64(N)

	// ---- Релионеаризация ----
	ctMul, _ = evaluator.MulNew(ct1, ct2)
	durRelin := time.Duration(0)
	ctRelin := bgv.NewCiphertext(params, 1, params.MaxLevel())
	for i := 0; i < N; i++ {
		start := time.Now()
		evaluator.Relinearize(ctMul, ctRelin)
		durRelin += time.Since(start)
	}
	grid["Релионеаризация"][key] = float64(durRelin.Microseconds()) / float64(N)

	// ---- Вращение (rotate rows 1 шаг) ----
	durRotate := time.Duration(0)
	ctRotate := bgv.NewCiphertext(params, 1, params.MaxLevel())
	for i := 0; i < N; i++ {
		start := time.Now()
		evaluator.RotateRows(ct1, ctRotate)
		durRotate += time.Since(start)
	}
	grid["Вращение (rotate rows 1 шаг)"][key] = float64(durRotate.Microseconds()) / float64(N)

	// ---- Сериализация шифротекста ----
	durSerialize := time.Duration(0)
	for i := 0; i < N; i++ {
		start := time.Now()
		ct1.MarshalBinary()
		durSerialize += time.Since(start)
	}
	grid["Сериализация шифротекста"][key] = float64(durSerialize.Microseconds()) / float64(N)

	// ---- Сжатая сериализация (ZLIB) ----
	durSerializeZLIB := time.Duration(0)
	for i := 0; i < N; i++ {
		data, _ := ct1.MarshalBinary()
		start := time.Now()
		compressZLIB(data)
		durSerializeZLIB += time.Since(start)
	}
	grid["Сжатая сериализация (ZLIB)"][key] = float64(durSerializeZLIB.Microseconds()) / float64(N)

	// ---- Сжатая сериализация (Zstandard) ----
	durSerializeZstd := time.Duration(0)
	for i := 0; i < N; i++ {
		data, _ := ct1.MarshalBinary()
		start := time.Now()
		compressZstd(data)
		durSerializeZstd += time.Since(start)
	}
	grid["Сжатая сериализация (Zstandard)"][key] = float64(durSerializeZstd.Microseconds()) / float64(N)
}

// benchCKKS выполняет бенчмарки операций схемы CKKS
func benchCKKS(params ckks.Parameters, key string, grid map[string]map[string]float64) {
	N := 100 // Количество итераций для усреднения

	// Создаем ключи один раз
	kgen := rlwe.NewKeyGenerator(params)
	sk := kgen.GenSecretKeyNew()

	// ---- Генерация релинеаризации ----
	durRelinKeyGen := time.Duration(0)
	var rlk *rlwe.RelinearizationKey
	for i := 0; i < N; i++ {
		start := time.Now()
		rlk = kgen.GenRelinearizationKeyNew(sk)
		durRelinKeyGen += time.Since(start)
	}
	grid["Генерация релинеаризации"][key] = float64(durRelinKeyGen.Microseconds()) / float64(N)

	// ---- Генерация Galois ключей ----
	durGaloisKeyGen := time.Duration(0)
	galEls := params.GaloisElementForComplexConjugation()
	for i := 0; i < N; i++ {
		start := time.Now()
		kgen.GenGaloisKeyNew(galEls, sk)
		durGaloisKeyGen += time.Since(start)
	}
	grid["Генерация Galois ключей"][key] = float64(durGaloisKeyGen.Microseconds()) / float64(N)

	// Создаем Galois ключ для вращений
	gk := kgen.GenGaloisKeyNew(galEls, sk)
	evk := rlwe.NewMemEvaluationKeySet(rlk, gk)

	// Подготовка объектов
	encoder := ckks.NewEncoder(params)
	encryptor := rlwe.NewEncryptor(params, sk)
	decryptor := rlwe.NewDecryptor(params, sk)
	evaluator := ckks.NewEvaluator(params, evk)

	// Создание тестовых данных
	values1 := make([]complex128, params.MaxSlots())
	values2 := make([]complex128, params.MaxSlots())
	r := rand.New(rand.NewSource(0))
	for i := range values1 {
		values1[i] = complex(2*r.Float64()-1, 2*r.Float64()-1)
		values2[i] = complex(2*r.Float64()-1, 2*r.Float64()-1)
	}

	pt1 := ckks.NewPlaintext(params, params.MaxLevel())
	pt2 := ckks.NewPlaintext(params, params.MaxLevel())
	encoder.Encode(values1, pt1)
	encoder.Encode(values2, pt2)

	// ---- Шифрование ----
	ct1, _ := encryptor.EncryptNew(pt1)
	ct2, _ := encryptor.EncryptNew(pt2)

	durEncrypt := time.Duration(0)
	for i := 0; i < N; i++ {
		start := time.Now()
		encryptor.EncryptNew(pt1)
		durEncrypt += time.Since(start)
	}
	grid["Шифрование"][key] = float64(durEncrypt.Microseconds()) / float64(N)

	// ---- Дешифрование ----
	durDecrypt := time.Duration(0)
	for i := 0; i < N; i++ {
		start := time.Now()
		decryptor.DecryptNew(ct1)
		durDecrypt += time.Since(start)
	}
	grid["Дешифрование"][key] = float64(durDecrypt.Microseconds()) / float64(N)

	// ---- Сложение ----
	durAdd := time.Duration(0)
	ctAdd := ckks.NewCiphertext(params, 1, params.MaxLevel())
	for i := 0; i < N; i++ {
		start := time.Now()
		evaluator.Add(ct1, ct2, ctAdd)
		durAdd += time.Since(start)
	}
	grid["Сложение"][key] = float64(durAdd.Microseconds()) / float64(N)

	// ---- Умножение ----
	durMul := time.Duration(0)
	ctMul := ckks.NewCiphertext(params, 2, params.MaxLevel())
	for i := 0; i < N; i++ {
		start := time.Now()
		evaluator.Mul(ct1, ct2, ctMul)
		durMul += time.Since(start)
	}
	grid["Умножение"][key] = float64(durMul.Microseconds()) / float64(N)

	// ---- Релионеаризация ----
	ctMul, _ = evaluator.MulNew(ct1, ct2)
	durRelin := time.Duration(0)
	ctRelin := ckks.NewCiphertext(params, 1, params.MaxLevel())
	for i := 0; i < N; i++ {
		start := time.Now()
		evaluator.Relinearize(ctMul, ctRelin)
		durRelin += time.Since(start)
	}
	grid["Релионеаризация"][key] = float64(durRelin.Microseconds()) / float64(N)

	// ---- Вращение (rotate rows 1 шаг) ----
	durRotate := time.Duration(0)
	ctRotate := ckks.NewCiphertext(params, 1, params.MaxLevel())
	for i := 0; i < N; i++ {
		start := time.Now()
		evaluator.Conjugate(ct1, ctRotate)
		durRotate += time.Since(start)
	}
	grid["Вращение (rotate rows 1 шаг)"][key] = float64(durRotate.Microseconds()) / float64(N)

	// ---- Сериализация шифротекста ----
	durSerialize := time.Duration(0)
	for i := 0; i < N; i++ {
		start := time.Now()
		ct1.MarshalBinary()
		durSerialize += time.Since(start)
	}
	grid["Сериализация шифротекста"][key] = float64(durSerialize.Microseconds()) / float64(N)

	// ---- Сжатая сериализация (ZLIB) ----
	durSerializeZLIB := time.Duration(0)
	for i := 0; i < N; i++ {
		data, _ := ct1.MarshalBinary()
		start := time.Now()
		compressZLIB(data)
		durSerializeZLIB += time.Since(start)
	}
	grid["Сжатая сериализация (ZLIB)"][key] = float64(durSerializeZLIB.Microseconds()) / float64(N)

	// ---- Сжатая сериализация (Zstandard) ----
	durSerializeZstd := time.Duration(0)
	for i := 0; i < N; i++ {
		data, _ := ct1.MarshalBinary()
		start := time.Now()
		compressZstd(data)
		durSerializeZstd += time.Since(start)
	}
	grid["Сжатая сериализация (Zstandard)"][key] = float64(durSerializeZstd.Microseconds()) / float64(N)
}
