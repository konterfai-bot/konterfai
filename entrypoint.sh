#!/bin/bash

/usr/local/bin/konterfai \
    --address="${ADDRESS:-0.0.0.0}" \
    --port="${PORT:-8080}" \
    --hallucinator-url=${HALLUCINATOR_URL:-"https://localhost:8080"} \
    --generate-interval="${GENERATE_INTERVAL:-2s}" \
    --hallucination-cache-size="${HALLUCINATION_CACHE_SIZE:-10}" \
    --hallucination-prompt-word-count="${HALLUCINATION_PROMPT_WORD_COUNT:-5}" \
    --hallucination-word-count="${HALLUCINATION_WORD_COUNT:-500}" \
    --hallucination-request-count="${HALLUCINATION_REQUEST_COUNT:-5}" \
    --hallucinator-link-percentage="${HALLUCINATOR_LINK_PERCENTAGE:-10}" \
    --hallucinator-link-max-subdirectory-depth="${HALLUCINATOR_LINK_MAX_SUBDIRECTORY_DEPTH:-5}" \
    --hallucinator-link-has-variables-probability="${HALLUCINATOR_LINK_HAS_VARIABLES_PROBABILITY:-0.5}" \
    --hallucinator-link-max-variables="${HALLUCINATOR_LINK_MAX_VARIABLES:-5}" \
    --ollama-address=${OLLAMA_ADDRESS:-"http://localhost:11434"} \
    --ollama-model=${OLLAMA_MODEL} \
    --ollama-request-timeout=${OLLAMA_REQUEST_TIMEOUT:-60s} \
    --ai-temperature=${AI_TEMPERATURE:-30.0} \
    --ai-seed=${AI_SEED:-0} \
    --webserver-200-probability=${WEBSERVER_200_PROBABILITY:-0.95} \
    --webserver-error-cache-size=${WEBSERVER_ERROR_CACHE_SIZE:-1000} \
    --random-uncertainty=${RANDOM_UNCERTAINTY:-0.1}
