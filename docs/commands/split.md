# ktidy split

Like `export`, but `--remove` also deletes the context from the source file.

```bash
ktidy split prod > prod.yaml                    # export only
ktidy split --remove prod > prod.yaml           # export and remove from source
ktidy split --remove --dry-run prod             # preview the removal
```
