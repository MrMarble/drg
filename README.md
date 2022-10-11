# Deep Rock Galactic
![DRG Version](https://img.shields.io/badge/DRG%20Version-1.36-yellow.svg?style=flat)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mrmarble/drg)

In progress tool to decompile save from game Deep Rock Galactic

# New version

I've moved from this CLI and built a web version, take a look at it! https://github.com/MrMarble/drg-editor


## Usage

I made a simple CLI tool to decompile save files.

```
Usage: DGR <file>

Manipulate Deep Rock Galactic save files

Arguments:
  <file>    Save file to manipulate

Flags:
  -h, --help       Show context-sensitive help.
      --meta       Print metadata
      --version    Print version information and quit
```

Example:

```shell
drg 76561198055109443_Player.sav > save.json
```

## UUIDS

Some of the values inside the save file are behind an UUID, this is the list of the known ones

### Classes

| UUID                             | Description |
|----------------------------------|-------------|
| 30d8ea17d8fbba4c95306de9655c2f8c | Scout       |
| 85ef626c65f1024a8dfeb5d0f3909d2e | Engineer    |
| 9edd56f1eebcc5488d5b5e5b80b62db4 | Driller     |
| ae56e180fec0c44d96fa29c28366b97b | Gunner      |

### Resources

| UUID                             | Description |
|----------------------------------|-------------|
| 078548b93232c04085f892e084a74100 | Yeast       |
| 22bc4f7d07d13e43bfca81bd9c14b1af | Jadiz       |
| 22daa757ad7a8049891b17edcc2fe098 | Barley      |
| 41ea550c1d46c54bbe2e9ca5a7accb06 | Malt        |
| 488d05146f5f754ba3d4610d08c0603e | Enor        |
| 5828652c9a5de845a9e2e1b8b463c516 | Error       |
| 5f2bcf8347760a42a23b6edc07c0941d | Umanita     |
| 67668aae828fdb48a9111e1b912dbfa4 | Phazyonite  |
| 72312204e287bc41815540a0cf881280 | Starch      |
| 8aa7fb43293a0b49b8be42ffe068a44c | Croppa      |
| 99fa526ad87748459498905a278693f6 | Data        |
| a10cb2853871fb499ac854a1cde2202c | Cores       |
| aaded8766c227d408032afd18d63561e | Magnite     |
| af0dc4fe8361bb48b32c92cc97e21de7 | Bismor      |

### Seasons

| UUID                             | Description |
|----------------------------------|-------------|
| a47d407ec0e4364892ce2e03de7df0b3 | Season 1    |
| b860b55f1d1bb54d8ee2e41fda9f5838 | Seaon 2     |

### Events

| UUID                             | Description       |
|----------------------------------|-------------------|
| af2717b41f16ba4faa2af86801fb522c | Space beach party |


### Schematics

| UUID                             | Description       |
|----------------------------------|-------------------|
| 5885a33b15ae844591a66b65a2e5494e | Chain Hit         |
| 8500b75544eeac41ad48486f36662352 |                   |
| 6f26a8b49f967c4d999f734645eae2c4 | Composite Casings |
| 4e86bcf8f790974cbb68f8ec7d9efde1 |                   |
