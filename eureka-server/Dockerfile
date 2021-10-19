FROM openjdk:8-jre-alpine

COPY app.jar /app.jar

ENV APP_OPTS=""
ENV JVM_OPTS="-Xss256k -XX:MaxRAMPercentage=80.0 -Duser.timezone=Asia/Shanghai -Djava.security.egd=file:/dev/./urandom"
ENV JAVA_OPTS=""
ENV EUREKA_INSTANCE_HOST_NAME=""
ENV EUREKA_INSTANCE_LIST=""

ENTRYPOINT ["sh", "-c", "java $JVM_OPTS  $JAVA_OPTS -jar /app.jar $APP_OPTS"]