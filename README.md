# Tribal Jokes API

Este es un proyecto de practica utilizando docker y go

#

### Archivo de configuración

En este proyecto se esta usando variables de entorno, la cuales se establecen en el archivo .env

1. PORT

   Esta variable especifica el puerto por el cual estará escuchando la aplicación

2. MAX_JOKES

   Esta variable establece el máximo de jokes que serán retornados

#

### Instrucciones

1. Clonar el repositorio

   ```
   git clone git@github.com:farodriguezm/tribal-jokes-api.git
   ```

2. Acceder a la directorio del proyecto

   ```
   cd tribal-jokes-api
   ```

3. Construir imagen de docker

   ```
   docker build --tag tribal-jokes-api .
   ```

4. Crear contenedor

   ```
   docker run -d -p 8000:8000 --env-file .env --name jokes-api tribal-jokes-api
   ```

#

### Rutas disponibles

1. /

   Retorna un mensaje de bienvenida

2. /ping

   Retorna el mensaje "pong"

3. /jokes/sync

   Retorna una lista de jokes de forma secuencial

4. /jokes/wg

   Retorna una lista de jokes utilizado WaitGroups

5. /jokes/chanel

   Retorna una lista de jokes utilizado Chanels
