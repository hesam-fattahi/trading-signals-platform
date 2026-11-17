# syntax=docker/dockerfile:1

# ==============================================================================
# 1 - Build the JAR
# ==============================================================================
FROM maven:3.9-eclipse-temurin-26 AS builder

WORKDIR /app

COPY proto/ /proto/

COPY services/indicator-calculator/pom.xml /app/pom.xml
RUN --mount=type=cache,target=/root/.m2 \
    mvn -B -q dependency:go-offline

COPY services/indicator-calculator/src/ /app/src/
RUN --mount=type=cache,target=/root/.m2 \
    mvn -B -q package -DskipTests

# ==============================================================================
# 2 - Runtime image
# ==============================================================================
FROM eclipse-temurin:25-jre-alpine

RUN apk add --no-cache --virtual .build-deps \
        gcc make musl-dev curl tar && \
    curl -fsSL https://downloads.sourceforge.net/project/ta-lib/ta-lib/0.4.0/ta-lib-0.4.0-src.tar.gz \
        -o /tmp/ta-lib.tar.gz && \
    tar -xzf /tmp/ta-lib.tar.gz -C /tmp && \
    cd /tmp/ta-lib && \
    ./configure --prefix=/usr && \
    make && make install && \
    apk del .build-deps && \
    rm -rf /tmp/ta-lib*

WORKDIR /app
COPY --from=builder /app/target/indicator-calculator-0.1.0.jar /app/app.jar

ENTRYPOINT ["java", "-jar", "/app/app.jar"]