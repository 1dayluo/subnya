<!--
 * @Author: 1dayluo
 * @Date: 2023-02-07 11:18:40
 * @LastEditTime: 2023-03-11 22:03:55
-->
# SubNya
## <div align="center"><b><a href="README.md">English</a> | <a href="README_CN.md">简体中文</a> | <a href="README_ES.md">Español</a></b></div>

![](https://img.shields.io/github/commit-activity/w/1dayluo/SubNya_monitor?style=flat-square)    ![](https://img.shields.io/github/license/1dayluo/SubNya_monitor?style=flat-square) 

## Introducción

SubNya_monitor es una nueva herramienta de enumeración y monitoreo de subdominios utilizada para rastrear el estado de los subdominios en el dominio objetivo, incluyendo subdominios recién añadidos y eliminados. Utiliza goroutines para aumentar la velocidad de enumeración de subdominios. La herramienta almacena datos tanto en Redis como en SQLite, aprovechando la velocidad de Redis para monitorear cambios en el MD5 de los archivos, y las características de SQLite para almacenar y actualizar los datos de subdominios, incluyendo el uso de transacciones. Finalmente, el resultado de la salida se guardará en un registro de archivo local (opcional) o se enviará como notificación a un Telegram/correo electrónico personal a través de una API.

El proyecto actual ha completado sus funcionalidades básicas; otras funciones y el Dockerfile aún están en desarrollo (ver la sección todo más abajo). 

## Instalación sencilla
### release

 Puedes descargarlos desde la página de [releases](https://github.com/1dayluo/SubNya_monitor/releases/tag/v1.0).

### usando go install 
Si ya tienes un entorno de Go listo (al menos go 1.19), es tan fácil como:
```lua
go install  github.com/1dayluo/subnya@latest  

```

## Archivo de Configuración

El archivo de configuración se encuentra en `~/.config/subnya/config/config.yml` y tiene este aspecto. Los archivos ejecutables deben incluir la ruta absoluta:

```yml
schedule:
  - cron: "* * * *"   # Para despliegue con Dockerfile

monitor:
  dir: 
    - "./test"  # La carpeta a monitorear - se recorrerán todos los archivos en esta carpeta
  settings:
    - timeout : 30   # Tiempo de intervalo
    - threads : 10    # Número de hilos
    - maxenumerationtime: 10   # Tiempo máximo de enumeración
    - outfile : "/var/tmp/"    # Carpeta de salida
  
  
redis:
  addr: "172.17.0.1:6379"   # Dirección de conexión de Redis
  password: ""   # Contraseña de Redis
  db:   0 

sqlite:
  db_1: "./db/monitor.db"   # Establecer la ubicación de almacenamiento de la base de datos
```

## Uso

```lua
Options:
  --update, -u           Verificar si el monitor tiene actualizaciones
  --run, -r              Iniciar el buscador de subdominios y actualizar datos (incluyendo código de estado de respuesta) en SQLite
  --output OUTPUT
  --help, -h             Mostrar esta ayuda y salir
```
