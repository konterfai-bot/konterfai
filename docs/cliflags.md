[<- back to docs](README.md)

<style>
    th {
        display: none;
    }
</style>

# Command Line Flags

The following flags are available for the konterfAI binary:

- `--address`

|                 |                                            |
|-----------------|--------------------------------------------|
| **Type:**       | string                                     |
| **Default:**    | `127.0.0.1`                                |
| **Description** | The address the webserver will run run on. |

- `--port`

|                |                                               |
|----------------|-----------------------------------------------|
| **Type:**      | integer                                       |
| **Default:**   | `8080`                                        |
| **Description**| The port konterfAIs webserver will listen on. |

- `--statistics-port`

|                |                                                                                                                                                 |
|----------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| **Type:**      | integer                                                                                                                                         |
| **Default:**   | `8081`                                                                                                                                          |
| **Description**| The port konterfAIs statistics webserver will listen on.<br> The `/metrics` endpoint for [prometheus](https://prometheus.io/) also lives there. |

- `--hallucinator-url`

|                 |                                                                                                                   |
|-----------------|-------------------------------------------------------------------------------------------------------------------|
| **Type:**       | url                                                                                                               |
| **Default:**    | http://localhost:8080                                                                                             |
| **Description** | The FQDN konterfAI uses. Must match the settings of your reverse proxy (if konterfAI is not running stand-alone). |

- `--generate-interval`

|                 |                                                                                                            |
|-----------------|------------------------------------------------------------------------------------------------------------|
| **Type:**       | integer                                                                                                    |
| **Default:**    | 5s                                                                                                         |
| **Description** | The interval in seconds to wait before attempting to generate a new hallucination, when the cache is full. |

- `--hallucination-cache-size`

|                  |                                                                                                                                               |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| **Type:**        | integer                                                                                                                                       |
| **Default:**     | 10                                                                                                                                            |
| **Description**  | The number of hallucinations to cache. Use high numbers for slow CPUs/GPUs and low numbers if you have vast amount of CPU-/GPU-time to spare. |

- `--hallucination-prompt-word-count`

|                 |                                                                                                                                                                                       |
|-----------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Type:**       | integer                                                                                                                                                                               |
| **Default:**    | 5                                                                                                                                                                                     |
| **Description** | The number of words (nouns, verbs, ..) to use for hallucination prompts. More words means a higher probabpility for the result to become a vivid hallucination (like a feaver-dream). |

- `--hallucination-word-count`

|                  |                                                                                                          |
|------------------|----------------------------------------------------------------------------------------------------------|
| **Type:**        | integer                                                                                                  |
| **Default:**     | 500                                                                                                      |
| **Description**  | The number of words that is expected from the resulting hallucination (length of the generated article). |

- `--hallucination-request-count`

|                  |                                                                                                                    |
|------------------|--------------------------------------------------------------------------------------------------------------------|
| **Type:**        | integer                                                                                                            |
| **Default:**     | 5                                                                                                                  |
| **Description**  | Counter how many times the same hallucination should be presented. Use a high number here to reduce CPU-/GPU-load. |

-- `--hallucination-minimal-length`

|                 |                                                                                    |
|-----------------|------------------------------------------------------------------------------------|
| **Type:**       | int                                                                                |
| **Default:**    | 500                                                                                |
| **Description** | The minimal length of a hallucination in characters. Use <1 to disable this check. |

- `--hallucinator-link-percentage`

|                  |                                                                              |
|------------------|------------------------------------------------------------------------------|
| **Type:**        | integer                                                                      |
| **Default:**     | 10                                                                           |
| **Description**  | The percentage of links to add to the hallucination measured by total words. |

- `--hallucinator-link-max-subdirectory-depth`

|                  |                                                                       |
|------------------|-----------------------------------------------------------------------|
| **Type:**        | integer                                                               |
| **Default:**     | 5                                                                     |
| **Description**  | The maximum number of subdirectories for a link in the hallucination. |

- `--hallucinator-link-has-variables-probability`

|                  |                                                         |
|------------------|---------------------------------------------------------|
| **Type:**        | float                                                   |
| **Default:**     | 0.5                                                     |
| **Description**  | The probability of a link having variables. (0.5 = 50%) |

- `--hallucinator-link-max-variables`

|                  |                                                                  |
|------------------|------------------------------------------------------------------|
| **Type:**        | integer                                                          |
| **Default:**     | 5                                                                |
| **Description**  | The maximum number of variables for a link in the hallucination. |

- `--ollama-address`

|                  |                                    |
|------------------|------------------------------------|
| **Type:**        | url                                |
| **Default:**     | http://localhost:11434             |
| **Description**  | The address of the ollama service. |

- `--ollama-model`

|                 |                                                                                          |
|-----------------|------------------------------------------------------------------------------------------|
| **Type:**       | string                                                                                   |
| **Default:**    | qwen2:0.5b                                                                               |
| **Description** | The model to use for hallucinations. Must be an active model in the ollama instance.<br>The smaller the model, the faster the hallucination generation will be and the less CPU-/GPU-time will be used. |

- `--ollama-request-timeout`

|                 |                                     |
|-----------------|-------------------------------------|
| **Type:**       | integer                             |
| **Default:**    | 60s                                 |
| **Description** | The timeout for the ollama service. |

- `--ai-temperature`

|                 |                                                                                                        |
|-----------------|--------------------------------------------------------------------------------------------------------|
| **Type:**       | float                                                                                                  |
| **Default:**    | 30.0                                                                                                   |
| **Description** | The temperature for the AI. Use a high number for more randomness and a low number for more coherence. |

- `--ai-seed`

|                 |                                   |
|-----------------|-----------------------------------|
| **Type:**       | int                               |
| **Default:**    | 0                                 |
| **Description** | The seed value to use for the AI. |

- `--webserver-200-probability`

|                 |                                                                            |
|-----------------|----------------------------------------------------------------------------|
| **Type:**       | float                                                                      |
| **Default:**    | 0.95                                                                       |
| **Description** | The probability of returning a 200 status code for a request (0.95 = 95%). |

- `--webserver-error-cache-size`

|                 |                                                                                                                                                                  |
|-----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Type:**       | integer                                                                                                                                                          |
| **Default:**    | 1000                                                                                                                                                             |
| **Description** | The number of error responses to cache (as long as an url is cached there, the request to that url would return the same error code if requested multiple times. |

- `--random-uncertainty`

|                 |                                                                                              |
|-----------------|----------------------------------------------------------------------------------------------|
| **Type:**       | float                                                                                        |
| **Default:**    | 0.1                                                                                          |
| **Description** | The uncertainty for the random generator (0.1 = 10%). Use a high number for more randomness. |

- `--tracing-endpoint`

|                 |                                                                                                                                                                        |
|-----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Type:**       | string                                                                                                                                                                 |
| **Default:**    |                                                                                                                                                                        |
| **Description** | The endpoint for the tracing server (open telemetry, e.g. [Jaeger](https://www.jaegertracing.io/)).<br/> If empty, tracing is disabled.<br/>Example value for Jaeger: --tracing-endpoint=localhost:4317 |

- `--log-level`

|                 |                                                                                   |
|-----------------|-----------------------------------------------------------------------------------|
| **Type:**       | string                                                                            |
| **Default:**    | info                                                                              |
| **Description** | The log level for the application. Possible values are: debug, info, warn, error. |

- `--log-format`

|                 |                                                                           |
|-----------------|---------------------------------------------------------------------------|
| **Type:**       | string                                                                    |
| **Default:**    | text                                                                      |
| **Description** | The log format for the application. Possible values are: json, text, off. |

<!-- Example table for easy copy & paste
|                 |   |
|-----------------|---|
| **Type:**       |   |
| **Default:**    |   |
| **Description** |   |
-->