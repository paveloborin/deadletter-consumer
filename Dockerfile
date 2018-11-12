FROM tutum/curl

# Environment Variables
ENV AMQP_HOST "rabbit"
ENV AMQP_PORT "5672"
ENV AMQP_USER "guest"
ENV AMQP_PASSWORD "guest"

ENV EXCHANGE_NAME "mailer"
ENV QUEUE_NAME "mailer"
ENV ROUTING_KEY ""

ENV DEAD_LETTER_EXCHANGE_NAME "mailer_fail"
ENV DEAD_LETTER_QUEUE_NAME "mailer_fail"

COPY bin/consumer /consumer

EXPOSE 8000

CMD ["/consumer", "--pretty-logging"]
