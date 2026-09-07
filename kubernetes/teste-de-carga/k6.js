import http from 'k6/http';
import { check } from 'k6';

const TEST_ID =
  __ENV.TEST_ID ||
  'projeto-korp';

export const options = {

  tags: {
    testid: TEST_ID,
  },

  stages: [

    {
      duration: '30s',
      target: 20,
    },

    {
      duration: '30s',
      target: 100,
    },

    {
      duration: '1m',
      target: 200,
    },

    {
      duration: '1m',
      target: 300,
    },

    {
      duration: '1m',
      target: 400,
    },

    {
      duration: '1m',
      target: 500,
    },

    {
      duration: '2m',
      target: 500,
    },

    {
      duration: '30s',
      target: 200,
    },

    {
      duration: '30s',
      target: 0,
    },

  ],

  thresholds: {

    http_req_failed: [
      'rate<0.05',
    ],

    http_req_duration: [
      'p(95)<500',
    ],

  },

};

const BASE_URL =
  __ENV.BASE_URL ||
  'http://192.168.56.21:30080';

export default function () {

  const response = http.get(
    `${BASE_URL}/projeto-korp`,
    {
      timeout: '5s',
    }
  );

  check(response, {

    'HTTP status 200': (r) =>
      r &&
      r.status === 200,

    'Resposta contem Projeto Korp': (r) =>
      r &&
      typeof r.body === 'string' &&
      r.body.includes('Projeto Korp'),

  });

}