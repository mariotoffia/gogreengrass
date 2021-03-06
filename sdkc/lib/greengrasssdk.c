/*
 * Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 */

#include <stdio.h>

#include "greengrasssdk.h"
#include <stdarg.h>
#include <string.h>

/*
 * In case the system loads the stub library instead of the true
 * implementation library shipped with GGC by mistake, print an error.
 */
static void print_loaded_stub_error() {
    gg_log(GG_LOG_ERROR, "Loaded stub instead of implementation library!\n");
}

/***************************************
**            Global Methods          **
***************************************/

gg_error gg_global_init(uint32_t opt) {
    (void)opt;
    print_loaded_stub_error();
    return GGE_RESERVED_MAX;
}

/***************************************
**           Logging Methods          **
***************************************/

gg_error gg_log(gg_log_level level, const char *format, ...) {

    va_list logArgs;
    va_start(logArgs, format);
    
    switch(level) {
        case GG_LOG_DEBUG:
            fprintf(stderr, "DEBUG ");
            break;
        case GG_LOG_INFO:
            fprintf(stderr, "INFO ");
            break;
        case GG_LOG_WARN:
            fprintf(stderr, "WARN ");
            break;
        case GG_LOG_ERROR:
            fprintf(stderr, "ERROR ");
            break;
        case GG_LOG_FATAL:
            fprintf(stderr, "FATAL ");
            break;
        default:
            break;
    }

    vfprintf(stderr, format, logArgs);

    va_end(logArgs);

    return GGE_RESERVED_MAX;
}

/***************************************
**         gg_request Methods         **
***************************************/

gg_error gg_request_init(gg_request *ggreq) {
    (void)ggreq;
    return GGE_SUCCESS;
}

gg_error gg_request_close(gg_request ggreq) {
    (void)ggreq;
    return GGE_SUCCESS;
}

gg_error gg_request_read(gg_request ggreq, void *buffer, size_t buffer_size,
                         size_t *amount_read) {
    (void)ggreq;
    (void)buffer;
    (void)buffer_size;
    (void)amount_read;
    print_loaded_stub_error();
    return GGE_RESERVED_MAX;
}

/***************************************
**           Runtime Methods          **
***************************************/

gg_error gg_runtime_start(gg_lambda_handler handler, uint32_t opt) {
    (void)handler;
    (void)opt;
    print_loaded_stub_error();
    gg_lambda_context ctx;

    ctx.client_context = "{\"custom\":{\"subject\":\"invoke/testlambda\"}}";
    ctx.function_arn = "arn:aws:lambda:eu-central-1:033549287452:function:testlambda:2";

    handler(&ctx);
    return GGE_SUCCESS;
}

int readcnt = 0;
char payload[] = "{\"data\":44, \"hello\":\"world\"}";

gg_error gg_lambda_handler_read(void *buffer, size_t buffer_size,
                                size_t *amount_read) {
    
    int left = strlen(payload) - readcnt;

    if (left == 0) {
        return GGE_SUCCESS;
    }

    char *ptr = payload;

    strncpy(buffer, (ptr + readcnt), left); 
    readcnt += left;
    *amount_read = left;

    return GGE_SUCCESS;
}

gg_error gg_lambda_handler_write_response(const void *response,
                                          size_t response_size) {
    (void)response;
    (void)response_size;

    gg_log(GG_LOG_INFO, "response: '%s'\n", (char*)response);

    return GGE_SUCCESS;
}

gg_error gg_lambda_handler_write_error(const char *error_message) {
    (void)error_message;
    
    gg_log(GG_LOG_WARN, "write error message: '%s'\n", error_message);

    return GGE_SUCCESS;
}

/***************************************
**     AWS Secrets Manager Methods    **
***************************************/

gg_error gg_get_secret_value(gg_request ggreq, const char *secret_id,
                             const char *version_id, const char *version_stage,
                             gg_request_result *result) {
    (void)ggreq;
    (void)secret_id;
    (void)version_id;
    (void)version_stage;
    (void)result;
    print_loaded_stub_error();
    return GGE_RESERVED_MAX;
}

/***************************************
**           Lambda Methods           **
***************************************/

gg_error gg_invoke(gg_request ggreq, const gg_invoke_options *opts,
                   gg_request_result *result) {
    (void)ggreq;
    (void)opts;
    (void)result;
    print_loaded_stub_error();
    return GGE_RESERVED_MAX;
}

/***************************************
**           AWS IoT Methods          **
***************************************/

gg_error gg_publish_options_init(gg_publish_options *opts) {
    (void)opts;
    return GGE_SUCCESS;
}

gg_error gg_publish_options_free(gg_publish_options opts) {
    (void)opts;
    print_loaded_stub_error();
    return GGE_SUCCESS;
}

gg_error gg_publish_options_set_queue_full_policy(gg_publish_options opts,
        gg_queue_full_policy_options policy) {
    (void)opts;
    (void)policy;
    return GGE_SUCCESS;
}

gg_error gg_publish_with_options(gg_request ggreq, const char *topic,
        const void *payload, size_t payload_size, const gg_publish_options opts,
        gg_request_result *result) {
    (void)ggreq;
    (void)topic;
    (void)payload;
    (void)payload_size;
    (void)opts;
    (void)result;
    return GGE_SUCCESS;
}

gg_error gg_publish(gg_request ggreq, const char *topic, const void *payload,
                    size_t payload_size, gg_request_result *result) {
    (void)ggreq;
    (void)topic;
    (void)payload;
    (void)payload_size;
    (void)result;
    print_loaded_stub_error();
    return GGE_RESERVED_MAX;
}

gg_error gg_get_thing_shadow(gg_request ggreq, const char *thing_name,
                             gg_request_result *result) {
    (void)ggreq;
    (void)thing_name;
    (void)result;
    return GGE_SUCCESS;
}

gg_error gg_update_thing_shadow(gg_request ggreq, const char *thing_name,
                                const char *update_payload,
                                gg_request_result *result) {
    (void)ggreq;
    (void)thing_name;
    (void)update_payload;
    (void)result;
    return GGE_SUCCESS;
}

gg_error gg_delete_thing_shadow(gg_request ggreq, const char *thing_name,
                                gg_request_result *result) {
    (void)ggreq;
    (void)thing_name;
    (void)result;
    return GGE_SUCCESS;
}
