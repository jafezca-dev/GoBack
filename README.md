# GoBack

Simple file backup program for Minio storage service.

| Params        | Description  |
| ------------- | ------------- |
| path          |   |
| client        |   |
| bucket        |   |
| endpoint      |   |
| accesskey     |   |
| secretkey     |   |
| ignorefolder  |   |
| ignorefile    |   |
| date          |   |
| uploadthread  |   |

### BackUp
```
goback (full/incremental) --client minio --bucket luminea --endpoint minio.goback.go --accesskey 49neh355PRK4AgQLT9f --secretkey 9CincmHVVvdRpiAkcfzy2bT421iLDU --path /home/user
```

### List BackUps
```
goback list --client minio --bucket luminea --endpoint minio.goback.go --accesskey 49neh355PRK4AgQLT9f --secretkey 9CincmHVVvdRpiAkcfzy2bT421iLDU
```

### Recovery
```
goback recovery --client minio --bucket luminea --endpoint minio.goback.go --accesskey 49neh355PRK4AgQLT9f --secretkey 9CincmHVVvdRpiAkcfzy2bT421iLDU --date 2025_01_08_03_09_05
```