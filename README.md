# tfs_item_attributes_reader

Unserializes/Serializes TFS database stored items.

You can directly load data as `[]byte`, load from hex string or load from a file.

## Example

```go
item, err := UnserializeHexString("18190073616D7572616920736F756C2072756E652028543129202B3822010000000000000005006C6576656C020800000000000000")
if err != nil {
        // Handle
        return
}

log.Println(item.Name)
...
```

The unserialized item will look like this:

```json
{
        "Name": "Super Secret Rune",
        "CustomAttributes": {
                "level": {
                        "Key": "level",
                        "Value": 8,
                }
        }
}
```

You can also serialize it back:

```go
s, err := item.SerializeAsHexString()
if err != nil {
        // Handle
        return
}

log.Println(s)
```

```txt
18190073616D7572616920736F756C2072756E652028543129202B3822010000000000000005006C6576656C020800000000000000
```
