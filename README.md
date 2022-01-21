# pxlmtc-test
Test interview for pixelmatic

### additional scope
- Acyclic graph ( *the value for each key will be the representation of a
NEW room* ).
- key is purely indicative, any value is tolerated (`forward`, `left`, `right`, `into the well`, `on the tree`, etc.).
- value is purely indicative except for `exit`, any value is tolerated.
- Due to input format, a value is a crossing (JSON object) OR a no-go zone (`deadend`, `dragon`, etc.), it can't be both.
- First exit found will be considered, other exits (even closer from entry) will be discarded.

### general instructions
This repository contains a maze `solver` and a maze `generator`.

`test/` directory contains multiple valid test files to test solver.

`test/input_00.json` to `test/input_02.json` are subject examples, other files were generated with `generator`.

`test/input_04.json` was manually modified to set exit at end of file to test worst case performance of solver.

### solver
```sh
> make tidy
> make solver
> cat test/input_00.json | ./bin/pxlmtc_solver
> cat test/input_04.json | ./bin/pxlmtc_solver # 90mb test file with exit in last lines
> cat test/input_04.json | ./bin/pxlmtc_solver -timeout=1s # timeout should raise on standard machine
```

### generator
```sh
> make tidy
> make generator
> ./bin/pxlmtc_generator -width=3 -height=5 # Due to random generation, you may have a small or empty maze, feel free to relaunch this command multiple times
> ./bin/pxlmtc_generator -width=3 -height=5 -exit=80 # exit will quickly appear (80%), "accidental" feature
> ./bin/pxlmtc_generator -width=5 -height=15 > test/input_XX.json # CAUTION !!! Huge file generation, it may vary due to randomness.
> ./bin/pxlmtc_generator -width=5 -height=20 -timeout=1s # timeout should raise on standard machine, it may vary due to randomness.
```
