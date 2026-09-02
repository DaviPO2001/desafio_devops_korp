import http from 'k6/http';
import { check } from 'k6';


// =========================================================
// CONFIGURACAO DO TESTE
// =========================================================

export const options = {

  stages: [

    // Inicio gradual
    {
      duration: '30s',
      target: 20,
    },

    // Aumenta a carga
    {
      duration: '1m',
      target: 100,
    },

    // Carga mais alta para tentar acionar o HPA
    {
      duration: '2m',
      target: 200,
    },

    // Reduz gradualmente
    {
      duration: '30s',
      target: 0,
    },

  ],


  // =======================================================
  // CRITERIOS DO TESTE
  // =======================================================

  thresholds: {

    // Menos de 1% das requisicoes podem falhar
    http_req_failed: [
      'rate<0.01',
    ],

    // 95% das requisicoes devem responder abaixo de 500ms
    http_req_duration: [
      'p(95)<500',
    ],

  },

};


// =========================================================
// ENDPOINT
// =========================================================

const BASE_URL =
  __ENV.BASE_URL ||
  'http://192.168.56.21:30080';


// =========================================================
// TESTE
// =========================================================

export default function () {

  const response = http.get(
    `${BASE_URL}/projeto-korp`
  );


  check(response, {

    'HTTP status 200': (r) =>
      r.status === 200,

    'Resposta contem Projeto Korp': (r) =>
      r.body.includes('Projeto Korp'),

  });

}