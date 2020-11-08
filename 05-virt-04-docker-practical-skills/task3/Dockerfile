FROM node

RUN git clone https://github.com/simplicitesoftware/nodejs-demo.git /app

WORKDIR /app

ENV VCAP_APP_PORT 3000
ENV VCAP_APP_HOST 0.0.0.0

EXPOSE ${VCAP_APP_PORT}

RUN npm install

CMD npm start
