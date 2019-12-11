# ERM-TOOLS 

## 使用配置

配置示例:

```json
{
    "ermFile": "*.erm",
    "newErmPath": "D:\\workspace\\JavaProjects\\demo\\src\\main\\resources\\db.erm",
    "oldErmPath": "D:\\workspace\\JavaProjects\\demo\\src\\main\\resources\\db2.erm",
    "dbName": "demodb",
    "dbHost": "",
    "dbUser": "",
    "dbPass": "",
    "dbType": "mysql",
    "type": "erm-erm",
    "outPath": "",
    "genDdl": true
}
```
    
说明:
- `ermFile` erm文件，可指定单个文件,*.erm 为所有文件
- `newErmPath` 当前erm文件路径
- `oldErmPath` 上一版本erm文件路径
- `dbName` DB名称
- `dbHost` DB host
- `dbUser` DB user
- `dbPass` DB password
- `dbType` mysql
- `type` erm-erm, erm-db 两种模式
- `outPath` 输出路径，默认当前执行路径
- `genDdl` true 生成 sql 文件

> type 为 erm-erm DB 相关可不配置，erm-db 时 oldErmPath 可不配置