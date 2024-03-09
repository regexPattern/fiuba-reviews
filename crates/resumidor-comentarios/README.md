# `resumidor-comentarios`

Utilidad para generar y actualizar los resúmenes de los comentarios de los docentes incrementalmente.

## Diseño

Cuando se inicia una nueva base de datos utilizando el [adaptador](https://github.com/regexPattern/fiuba-reviews?tab=readme-ov-file#adaptaci%C3%B3n-de-los-datos-originales), ningún docente cuenta con su resumen de comentarios ya generado, debido a que esta información no se obtiene directamente del scraping que se hace a los datos originales de Dolly.

Por lo tanto, los resúmenes deben ser generados una vez se hayan obtenido los comentarios de todos los docentes. Sin embargo, este proceso no puede ser automatizado al momento de popular la base de datos inicial porque no es un proceso que se pueda planificar de manera programática, debido a posibles errores que puedan ocurrir durante la generación de los resúmenes que dejarían el SQL de inicialización de la base de datos en un estado incompleto o inválido. Por ejemplo, podría ocurrir que no se completen los resúmenes de todos los docentes debido a algún rate-limiter que tenga la API del modelo usado.

Debido a esto, se decidió que la utilidad se pudiera ejecutar a conveniencia en lugar de hacerse automáticamente. La aplicación actualiza los resúmenes de la mayor cantidad de docentes posible de una vez, pero en caso de que ocurran errores y se interrumpa el proceso de actualización, no se cancela el proceso entero. En cambio, se actualizan efectivamente los resúmenes de los docentes que se hayan generado. La próxima vez que se ejecute, el programa continuará con los docentes que quedaron pendientes.

## Utilización

Para actualizar la base de datos se necesita establecer conexión con la misma, por lo que es necesario establecer la variable de entorno `DATABASE_URL` antes de ejecutar el programa.

La configuración para utilizar el modelo generador de resúmenes, depende del modelo elegido.

Los dos modelos ya implementados son el modelo [BART](https://huggingface.co/facebook/bart-large-cnn) de Facebook y [GPT-3.5 Turbo](https://platform.openai.com/docs/models/gpt-3-5-turbo) de OpenAI. Inicialmente la aplicación utilizó el primero de estos modelos, ya que es open source y de uso gratuito a través del servicio [Inference API](https://huggingface.co/docs/api-inference/en/index) de Hugging Face. La versión más reciente utiliza el segundo, a través de [OpenAI API](https://platform.openai.com/docs/overview), ya que, a pesar de que es un modelo de pago, personalmente considero que la calidad de los resúmenes generados por el mismo justifica la inversión, aunque por esto no quiero decir que el modelo gratuito haya hecho un mal trabajo.

### OpenAI GPT-3.5 Turbo

Este modelo es el que se utiliza por defecto. Basta con definir la variable de entorno `OPENAI_API_KEY`, cuyo valor debe ser la clave de acceso a la OpenAI API. Para conseguir una clave, revisa la [documentación de OpenAI API](https://platform.openai.com/api-keys).

Podes correr la aplicación modificando el siguiente comando y reemplazando las variables correspondientes:

```bash
export DATABASE_URL=...
export OPENAI_API_KEY=...

cargo run --release
```

### Facebook BART

Para reemplazar el modelo por defecto por este segundo modelo implementado, hay que modificar el archivo que define el crate binario del programa ([`main.rs`](https://github.com/regexPattern/fiuba-reviews/tree/main/crates/resumidor-comentarios/src/gpt/hugging_face.rs)). Basta con reemplazar el cliente de OpenAI (`OpenAIClient`) que utiliza la versión actual del archivo, con el cliente de Inference API (`HuggingFaceClient`).

Para inicializar este modelo, se requiere la clave de acceso de la Inference API de Hugging Face. Podes conseguir dicha clave siguiendo la [documentación de Inference API](https://huggingface.co/docs/api-inference/en/quicktour#get-your-api-token).

Este es un ejemplo de cómo configurarla a través de una variable de entorno llamada `INFERENCE_API_KEY`:

```rust
use resumidor_comentarios::gpt::HuggingFaceClient;

let modelo = HuggingFaceClient {
    api_key: env::var("INFERENCE_API_KEY").unwrap(),
};
```

## Configuración

### Agregar nuevos modelos

Si querés agregar un nuevo modelo, podes agregarlo en un archivo dentro de la carpeta [`src/gpt`](https://github.com/regexPattern/fiuba-reviews/tree/main/crates/resumidor-comentarios/src/gpt/), que es donde están definidos los dos modelos ya implementados. Lo único que se necesita para crear un nuevo cliente del modelo a elección es implementar el trait `Modelo` definido en el archivo [`gpt.rs`](https://github.com/regexPattern/fiuba-reviews/tree/main/crates/resumidor-comentarios/src/gpt.rs).

Este trait requiere de la definición de una única función que toma una referencia a un cliente HTTP como parámetro, ya que por el momento solo se contemplaron modelos de IA que sean accedidos mediante el llamado a una API web externa.

### Condiciones de actualización de docentes

Para no estar teniendo que regenerar el resumen de comentarios de un docente cada vez que se le agrega un comentario, solo se someten a actualización los docentes que cumplen cierta condición.

Dicha condición consiste en un número mínimo de comentarios para generar el primer resumen de comentarios del docente, si este aún no tiene uno generado, y también una proporción mínima para regeneraciones subsecuentes del resumen. Es decir, si, por ejemplo, esta proporción se establece en `2/1`, el docente será considerado para actualización solo si actualmente tiene más del doble de comentarios de los que tenía cuando recibió su última actualización.

Ambas constantes son definidas en el archivo [`lib.rs`](https://github.com/regexPattern/fiuba-reviews/tree/main/crates/resumidor-comentarios/src/lib.rs) con los siguientes valores:

```rust
const MIN_COMENTARIOS_ACTUALIZACION: usize = 3;
const PROPORCION_COMENTARIOS_ACTUALIZACION: usize = 2 / 1;
```

### Máximo de solicitudes concurrentes

La API externa utilizada para llamar al modelo generativo puede tener restricciones en cuando al número de solicitudes concurrentes. Para solventar esto la constante `MAX_SOLICITUDES_CONCURRENTES` (también definida en el archivo [`lib.rs`](https://github.com/regexPattern/fiuba-reviews/tree/main/crates/resumidor-comentarios/src/lib.rs)) establece el límite a respetar.
