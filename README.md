# caserunner

## Compilación

```bash
go build -o caserunner ./cmd/caserunner
```

## Uso

```
caserunner [flags] <archivo-de-prueba> <codigo>
```

**Argumentos:**

| Argumento | Descripción |
|---|---|
| `<archivo-de-prueba>` | Ruta al archivo de prueba |
| `<codigo>` | Ruta al archivo de código que se va a probar |

**Flags:**

| Flag | Descripción |
|---|---|
| `-verbose` | Muestra el input y output completo de cada prueba |
| `-halt` | Detiene la ejecución al encontrar el primer error |

**Ejemplo:**

```bash
caserunner examples/fizzbuzz/test examples/fizzbuzz/fizzbuzz.py
```

## Formato del archivo de prueba

El archivo de prueba tiene dos partes: la **configuración** y los **casos de prueba**.

### Configuración

```
exec: python $code
time-limit: 1s
```

| Directiva | Requerida | Descripción |
|---|---|---|
| `exec:` | Sí | Comando para ejecutar el programa. `$code` se reemplaza con la ruta al archivo de código. |
| `time-limit:` | No | Tiempo máximo de respuesta por línea de output. Acepta unidades como `1s`, `500ms`, `3ms`. Si se omite, no hay límite. |

El texto fuera de los bloques de prueba se ignora, por lo que puedes usarlo como comentarios.

### Casos de prueba

Cada caso de prueba va entre delimitadores `--`:

```
--
input:
<lineas de entrada>

output:
<lineas de salida esperadas>
--
```

- Las secciones `input:` y `output:` pueden ir en cualquier orden dentro del bloque.
- Las líneas en blanco dentro de las secciones se ignoran.
- El orden de `input:` y `output:` importa: la N-ésima línea de input corresponde a la N-ésima línea de output.

### Líneas ignoradas en el output: `\`

Cuando una línea de input no produce output (por ejemplo, el número de casos de prueba al inicio de un problema), escribe `\` en la posición correspondiente dentro de `output:`.

Esto le indica al runner que no debe esperar respuesta del programa para esa línea de input, y tampoco la incluye al comparar el resultado.

```
--
input:
5        <-- cantidad de casos, no produce output
2 5 9 6 1
0

output:
\        <-- ignorar, no se espera output para "5"
\        <-- ignorar, no se espera output para "2 5 9 6 1"
2 6      <-- output esperado para "0"
--
```

### Escapar tokens especiales

Si una línea de output contiene literalmente `input:` u `output:`, escápala con `\` al inicio:

```
--
input:
3

output:
\output: Fizz
--
```

Esto produce el output esperado `output: Fizz` sin que el parser lo interprete como una nueva sección.
