# Tools

## clean_data

Cleans `world_measurements.json` and `world_durations.json` by removing items with obscure measurements, unit-derived names, and duplicates. Also renames unit-specific items to generic versions.

```bash
python3 tools/clean_data/clean_data.py
```

## mark_proper_nouns

Marks proper nouns in `world_measurements.json` and `world_durations.json` by adding `"proper_noun": true` to named entities (e.g., "Eiffel Tower", "Apollo 11"). Uses category-based rules and an explicit name set. Run after `clean_data`.

```bash
python3 tools/clean_data/mark_proper_nouns.py
```
