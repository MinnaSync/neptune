FROM oven/bun:latest

RUN mkdir -p /neptune
WORKDIR /neptune
COPY . /neptune/

RUN bun install

CMD ["bun", "run", "./src/index.ts"]