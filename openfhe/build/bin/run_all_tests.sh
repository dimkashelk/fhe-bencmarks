#!/bin/zsh
mkdir -p logs

echo "Starting OpenFHE benchmarks..."
cd benchmark

for f in *; do
  if [[ -x "$f" ]]; then
    echo ">>> Running $f"
    ./$f > "../logs/${f}.txt" 2>&1
    echo "Saved results to logs/${f}.txt"
  fi
done

echo "✅ All benchmarks completed. Results saved in ./logs/"
