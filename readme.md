# Hello-World Project for Watermill

This is a Hello-World project to study how to use Watermill as pubsub library.

* in this evaluation, will use go channel as pubsub adapter

## Adaptor

With gochannel adaptor:
* OutputChannelBuffer can set the output channel size; however, during graceful shutdown, not all the buffered messages are processed
    + increase OutputChannelBuffer can somehow increase buffered message, but not all
    + e.g. when OutputChannelBuffer is set to 20 and publish 20 messages, Ctrl+C when processing up to 3rd message, total 8 pending messages are processed before shutdown while other 9 messages are missed
* found that graceful shutdown might leave some pending messages not processed
* resuming after abrupt shutdown cannot resume processing the missed message

With redis adaptor:
* found that resuming after Ctrl+C shutdown might have last sent message resend
    + e.g. publish 20 messages, Ctrl+C when processing up to 4th message, graceful shutdown after processed 4 message; resuming the service successfully continue process the rest of the messages but have the 3rd and 4th message process again which was already processed before resume
* resuming after abrupt shutdown can resume processing the missed message with following problems
    + last processed message(s) might be process again
    + the order of the messages is changed
* here is one of the example log
    ```shell
    WKHKNB0018:hello-watermill kamching$ go run *.go --driver redis router
    [watermill] 2024/04/24 17:23:51.700188 router.go:264:   level=INFO  msg="Adding handler" handler_name=wk.email.send_handler topic=wk.email.send 
    [watermill] 2024/04/24 17:23:51.701745 router.go:264:   level=INFO  msg="Adding handler" handler_name=wk.imos.post_handler topic=wk.imos.post 
    [watermill] 2024/04/24 17:23:51.701964 router.go:398:   level=INFO  msg="Running router handlers" count=2 
    [watermill] 2024/04/24 17:23:51.702019 subscriber.go:182:       level=INFO  msg="Subscribing to redis stream topic" consumer_group=wk_api_consumer_group consumer_uuid=uWRa2UgfdPRxscPWaR7yEg provider=redis topic=wk.email.send 
    [watermill] 2024/04/24 17:23:51.702052 subscriber.go:203:       level=INFO  msg="Starting consuming" consumer_group=wk_api_consumer_group consumer_uuid=uWRa2UgfdPRxscPWaR7yEg provider=redis topic=wk.email.send 
    [watermill] 2024/04/24 17:23:51.715172 subscriber.go:182:       level=INFO  msg="Subscribing to redis stream topic" consumer_group=wk_api_consumer_group consumer_uuid=uWRa2UgfdPRxscPWaR7yEg provider=redis topic=wk.imos.post 
    [watermill] 2024/04/24 17:23:51.715214 subscriber.go:203:       level=INFO  msg="Starting consuming" consumer_group=wk_api_consumer_group consumer_uuid=uWRa2UgfdPRxscPWaR7yEg provider=redis topic=wk.imos.post 
    [watermill] 2024/04/24 17:23:51.715475 router.go:598:   level=INFO  msg="Starting handler" subscriber_name=wk.email.send_handler topic=wk.email.send 
    [watermill] 2024/04/24 17:23:51.721592 router.go:598:   level=INFO  msg="Starting handler" subscriber_name=wk.imos.post_handler topic=wk.imos.post 
    2024/04/24 17:23:51 Publish 0 to wk.email.send topic 981d1696-fc5a-41e4-a1bc-3102d221c62e
    2024/04/24 17:23:51 Send email: 981d1696-fc5a-41e4-a1bc-3102d221c62e message: Hello, Wah Kwong! (0)
    2024/04/24 17:23:51 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (0)
    2024/04/24 17:23:52 Publish 1 to wk.email.send topic 68841716-86ea-4647-b5c3-d93ca1cd8019
    2024/04/24 17:23:52 Send email: 68841716-86ea-4647-b5c3-d93ca1cd8019 message: Hello, Wah Kwong! (1)
    2024/04/24 17:23:52 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (1)
    2024/04/24 17:23:52 Publish 2 to wk.email.send topic 25f13fec-e359-4853-a5fc-f81e4b623b47
    2024/04/24 17:23:52 Send email: 25f13fec-e359-4853-a5fc-f81e4b623b47 message: Hello, Wah Kwong! (2)
    2024/04/24 17:23:52 Publish 3 to wk.email.send topic cf3363c4-3ed6-419a-b5f8-44f014ef98db
    2024/04/24 17:23:52 Publish 4 to wk.email.send topic a04f869a-8073-4e8b-82d8-29630b5f866c
    2024/04/24 17:23:52 Publish 5 to wk.email.send topic 6ca73411-5482-49d8-956f-0d8a20054f41
    2024/04/24 17:23:53 Publish 6 to wk.email.send topic f2f6b7f3-6314-420d-a1f2-07239ef0ac19
    2024/04/24 17:23:53 Publish 7 to wk.email.send topic e6600b2b-de4c-4b5b-bd6f-ab2449680c99
    2024/04/24 17:23:53 Publish 8 to wk.email.send topic ae38b88f-d2c3-41ac-bd73-f59b9df1beb1
    2024/04/24 17:23:53 Publish 9 to wk.email.send topic ce4d9488-2f09-4bbc-a89a-c1ebe916ac58
    2024/04/24 17:23:53 Publish 10 to wk.email.send topic 213ef9df-8a62-4851-aaa0-4835c4ba90b2
    2024/04/24 17:23:54 Publish 11 to wk.email.send topic 5b410366-fd7f-4eec-8ab1-34ce5fa83c2c
    2024/04/24 17:23:54 Publish 12 to wk.email.send topic 8648948f-4952-4abb-8f6a-f473623037dd
    2024/04/24 17:23:54 Publish 13 to wk.email.send topic 551ccf23-11bc-48e1-838c-41b57a484412
    2024/04/24 17:23:54 Publish 14 to wk.email.send topic 5e58d07e-5bde-49a4-b490-9f7c49a2b418
    2024/04/24 17:23:54 Publish 15 to wk.email.send topic 62222504-751e-45d4-bfb7-68e6c140b704
    2024/04/24 17:23:55 Publish 16 to wk.email.send topic fb60b399-4892-446d-a626-f3ad17e26809
    2024/04/24 17:23:55 Publish 17 to wk.email.send topic 43c89c51-c46b-43f6-a648-8c0498816420
    2024/04/24 17:23:55 Publish 18 to wk.email.send topic 0647924e-bc9e-48ee-8f92-c8d4ad1b062d
    2024/04/24 17:23:55 Publish 19 to wk.email.send topic 0f6dde2a-e8cc-455e-9ea3-a00f094a3edb
    ^C[watermill] 2024/04/24 17:23:56.383427 signals.go:19:         level=INFO  msg="Received interrupt signal, closing
    " 
    [watermill] 2024/04/24 17:23:56.383462 router.go:530:   level=INFO  msg="Closing router" 
    [watermill] 2024/04/24 17:23:56.383712 router.go:441:   level=INFO  msg="Subscriber stopped" subscriber_name=wk.email.send_handler topic=wk.email.send 
    [watermill] 2024/04/24 17:23:56.383732 router.go:441:   level=INFO  msg="Subscriber stopped" subscriber_name=wk.imos.post_handler topic=wk.imos.post 
    2024/04/24 17:23:56 Send email: cf3363c4-3ed6-419a-b5f8-44f014ef98db message: Hello, Wah Kwong! (3)
    [watermill] 2024/04/24 17:23:56.383700 router.go:377:   level=INFO  msg="Waiting for messages" timeout=30s 
    2024/04/24 17:23:56 Done sending email: 981d1696-fc5a-41e4-a1bc-3102d221c62e
    2024/04/24 17:23:56 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (2)
    [watermill] 2024/04/24 17:23:56.885309 subscriber.go:340:       level=ERROR msg="read fail" consumer_group=wk_api_consumer_group consumer_uuid=uWRa2UgfdPRxscPWaR7yEg err="context canceled" provider=redis topic=wk.imos.post 
    2024/04/24 17:23:57 Done sending email: 68841716-86ea-4647-b5c3-d93ca1cd8019
    2024/04/24 17:23:57 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (3)
    [watermill] 2024/04/24 17:23:57.018001 router.go:383:   level=INFO  msg="All messages processed" 
    [watermill] 2024/04/24 17:23:57.018071 router.go:541:   level=INFO  msg="Router closed" 
    2024/04/24 17:24:01 Done sending email: 25f13fec-e359-4853-a5fc-f81e4b623b47
    2024/04/24 17:24:02 Done sending email: cf3363c4-3ed6-419a-b5f8-44f014ef98db



    WKHKNB0018:hello-watermill kamching$ go run *.go --driver redis router
    [watermill] 2024/04/24 17:24:13.404963 router.go:264:   level=INFO  msg="Adding handler" handler_name=wk.email.send_handler topic=wk.email.send 
    [watermill] 2024/04/24 17:24:13.406406 router.go:264:   level=INFO  msg="Adding handler" handler_name=wk.imos.post_handler topic=wk.imos.post 
    [watermill] 2024/04/24 17:24:13.406597 router.go:398:   level=INFO  msg="Running router handlers" count=2 
    [watermill] 2024/04/24 17:24:13.406614 subscriber.go:182:       level=INFO  msg="Subscribing to redis stream topic" consumer_group=wk_api_consumer_group consumer_uuid=pjyAFSfJ2MigxmR89UdwLA provider=redis topic=wk.email.send 
    [watermill] 2024/04/24 17:24:13.406644 subscriber.go:203:       level=INFO  msg="Starting consuming" consumer_group=wk_api_consumer_group consumer_uuid=pjyAFSfJ2MigxmR89UdwLA provider=redis topic=wk.email.send 
    [watermill] 2024/04/24 17:24:13.414436 subscriber.go:182:       level=INFO  msg="Subscribing to redis stream topic" consumer_group=wk_api_consumer_group consumer_uuid=pjyAFSfJ2MigxmR89UdwLA provider=redis topic=wk.imos.post 
    [watermill] 2024/04/24 17:24:13.414490 subscriber.go:203:       level=INFO  msg="Starting consuming" consumer_group=wk_api_consumer_group consumer_uuid=pjyAFSfJ2MigxmR89UdwLA provider=redis topic=wk.imos.post 
    [watermill] 2024/04/24 17:24:13.414560 router.go:598:   level=INFO  msg="Starting handler" subscriber_name=wk.email.send_handler topic=wk.email.send 
    [watermill] 2024/04/24 17:24:13.415409 router.go:598:   level=INFO  msg="Starting handler" subscriber_name=wk.imos.post_handler topic=wk.imos.post 
    2024/04/24 17:24:13 Send email: 6ca73411-5482-49d8-956f-0d8a20054f41 message: Hello, Wah Kwong! (5)
    2024/04/24 17:24:13 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (5)
    2024/04/24 17:24:13 Send email: f2f6b7f3-6314-420d-a1f2-07239ef0ac19 message: Hello, Wah Kwong! (6)
    2024/04/24 17:24:13 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (6)
    2024/04/24 17:24:13 Send email: e6600b2b-de4c-4b5b-bd6f-ab2449680c99 message: Hello, Wah Kwong! (7)
    2024/04/24 17:24:18 Done sending email: 6ca73411-5482-49d8-956f-0d8a20054f41
    2024/04/24 17:24:18 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (7)
    2024/04/24 17:24:18 Send email: ae38b88f-d2c3-41ac-bd73-f59b9df1beb1 message: Hello, Wah Kwong! (8)
    2024/04/24 17:24:18 Done sending email: f2f6b7f3-6314-420d-a1f2-07239ef0ac19
    2024/04/24 17:24:18 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (8)
    2024/04/24 17:24:18 Send email: ce4d9488-2f09-4bbc-a89a-c1ebe916ac58 message: Hello, Wah Kwong! (9)
    2024/04/24 17:24:23 Done sending email: e6600b2b-de4c-4b5b-bd6f-ab2449680c99
    2024/04/24 17:24:23 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (9)
    2024/04/24 17:24:23 Done sending email: ae38b88f-d2c3-41ac-bd73-f59b9df1beb1
    2024/04/24 17:24:23 Send email: 213ef9df-8a62-4851-aaa0-4835c4ba90b2 message: Hello, Wah Kwong! (10)
    2024/04/24 17:24:23 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (10)
    2024/04/24 17:24:23 Send email: 5b410366-fd7f-4eec-8ab1-34ce5fa83c2c message: Hello, Wah Kwong! (11)
    2024/04/24 17:24:28 Done sending email: ce4d9488-2f09-4bbc-a89a-c1ebe916ac58
    2024/04/24 17:24:28 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (11)
    2024/04/24 17:24:28 Send email: 8648948f-4952-4abb-8f6a-f473623037dd message: Hello, Wah Kwong! (12)
    2024/04/24 17:24:28 Done sending email: 213ef9df-8a62-4851-aaa0-4835c4ba90b2
    2024/04/24 17:24:28 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (12)
    2024/04/24 17:24:28 Send email: 551ccf23-11bc-48e1-838c-41b57a484412 message: Hello, Wah Kwong! (13)
    2024/04/24 17:24:33 Done sending email: 5b410366-fd7f-4eec-8ab1-34ce5fa83c2c
    2024/04/24 17:24:33 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (13)
    2024/04/24 17:24:33 Send email: 5e58d07e-5bde-49a4-b490-9f7c49a2b418 message: Hello, Wah Kwong! (14)
    2024/04/24 17:24:33 Done sending email: 8648948f-4952-4abb-8f6a-f473623037dd
    2024/04/24 17:24:33 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (14)
    2024/04/24 17:24:33 Send email: 62222504-751e-45d4-bfb7-68e6c140b704 message: Hello, Wah Kwong! (15)
    2024/04/24 17:24:38 Done sending email: 551ccf23-11bc-48e1-838c-41b57a484412
    2024/04/24 17:24:38 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (15)
    2024/04/24 17:24:38 Send email: fb60b399-4892-446d-a626-f3ad17e26809 message: Hello, Wah Kwong! (16)
    2024/04/24 17:24:38 Done sending email: 5e58d07e-5bde-49a4-b490-9f7c49a2b418
    2024/04/24 17:24:38 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (16)
    2024/04/24 17:24:38 Send email: 43c89c51-c46b-43f6-a648-8c0498816420 message: Hello, Wah Kwong! (17)
    2024/04/24 17:24:43 Done sending email: 62222504-751e-45d4-bfb7-68e6c140b704
    2024/04/24 17:24:43 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (17)
    2024/04/24 17:24:43 Send email: 0647924e-bc9e-48ee-8f92-c8d4ad1b062d message: Hello, Wah Kwong! (18)
    2024/04/24 17:24:43 Done sending email: fb60b399-4892-446d-a626-f3ad17e26809
    2024/04/24 17:24:43 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (18)
    2024/04/24 17:24:43 Send email: 0f6dde2a-e8cc-455e-9ea3-a00f094a3edb message: Hello, Wah Kwong! (19)
    2024/04/24 17:24:48 Done sending email: 43c89c51-c46b-43f6-a648-8c0498816420
    2024/04/24 17:24:48 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (19)
    2024/04/24 17:24:48 Done sending email: 0647924e-bc9e-48ee-8f92-c8d4ad1b062d
    2024/04/24 17:24:53 Send email: 25f13fec-e359-4853-a5fc-f81e4b623b47 message: Hello, Wah Kwong! (2)
    2024/04/24 17:24:53 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (2)
    2024/04/24 17:24:53 Send email: cf3363c4-3ed6-419a-b5f8-44f014ef98db message: Hello, Wah Kwong! (3)
    2024/04/24 17:24:53 Done sending email: 0f6dde2a-e8cc-455e-9ea3-a00f094a3edb
    2024/04/24 17:24:53 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (3)
    2024/04/24 17:24:53 Send email: a04f869a-8073-4e8b-82d8-29630b5f866c message: Hello, Wah Kwong! (4)
    2024/04/24 17:24:58 Done sending email: 25f13fec-e359-4853-a5fc-f81e4b623b47
    2024/04/24 17:24:58 Send email to: chingkamhing@gmail.com body: Hello, Wah Kwong! (4)
    2024/04/24 17:24:58 Done sending email: cf3363c4-3ed6-419a-b5f8-44f014ef98db
    2024/04/24 17:25:03 Done sending email: a04f869a-8073-4e8b-82d8-29630b5f866c
    ```

## References

* [Getting started](https://watermill.io/docs/getting-started/)
