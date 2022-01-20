# pxlmtc-test
Test interview for pixelmatic

# assumptions
- Acyclic graph ( **the value for each key will be the representation of a
NEW room** ).
- key is purely indicative, any value is tolerated (`forward`, `left`, `right`, `into the well`, `on the tree`, etc.).
- value is purely indicative except for `exit`, any value is tolerated.
- Due to input format, a value is a crossing (JSON object) OR a no-go zone (`deadend`, `dragon`, etc.), it can't be both.
- First exit found will be considered, other exits (even closer from entry) will be discarded.
