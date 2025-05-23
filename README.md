# Series Tracker Backend (WEB_Lab6_Backend)

Este proyecto implementa una API RESTful en Go (Golang) para gestionar una colección de series de TV. Permite realizar operaciones CRUD (Crear, Leer, Actualizar, Eliminar) sobre las series, así como actualizar su estado, marcar episodios vistos y votar por ellas.

Utiliza PostgreSQL como base de datos y Docker Compose para facilitar el despliegue y la configuración del entorno.

## Características

* **Gestión de Series:**
    * Crear nuevas series.
    * Obtener la lista completa de series.
    * Obtener los detalles de una serie específica por ID.
    * Actualizar la información completa de una serie.
    * Eliminar una serie.
* **Acciones Específicas:**
    * Actualizar el estado de visualización de una serie (`Plan to Watch`, `Watching`, `Completed`, `Dropped`).
    * Incrementar el contador del último episodio visto.
    * Incrementar (Upvote) el ranking de una serie.
    * Decrementar (Downvote) el ranking de una serie.
* **API RESTful:** Endpoints claros y siguiendo las convenciones REST.
* **Persistencia:** Almacenamiento de datos en una base de datos PostgreSQL.
* **Contenerización:** Configuración de Docker y Docker Compose para un fácil despliegue.
* **CORS Habilitado:** Permite peticiones desde cualquier origen (configurable en `main.go`).

## Coleccion de postman
https://www.postman.com/iris-6065823/traker-series-iris/collection/pyks1gj/tracker-series?action=share&creator=43600137 

## Instalación y Ejecución

1.  **Clona el repositorio (si aplica):**
    ```bash
    git clone <URL_DEL_REPOSITORIO>
    cd WEB_Lab6_Backend
    ```
    O asegúrate de estar en el directorio raíz del proyecto (`WEB_Lab6_Backend`).

2.  **Construye e inicia los contenedores:**
    Este comando construirá la imagen Docker para el backend (si no existe o si el Dockerfile ha cambiado), creará e iniciará los contenedores para el backend y la base de datos PostgreSQL.
    ```bash
    docker-compose up --build -d
    ```
    * `--build`: Fuerza la reconstrucción de la imagen del backend.
    * `-d`: Ejecuta los contenedores en segundo plano (detached mode).

3.  **Para detener los servicios:**
    ```bash
    docker-compose down
    ```
    Si también quieres eliminar los volúmenes (¡perderás los datos de la BD!):
    ```bash
    docker-compose down -v
    ```
