[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/mRc53rxV)

en funcionamiento:
https://youtu.be/zXQKFqheFH0

Example_output.txt contiene un ejemplo del output entregado por el programa
debido a la concurrencia a veces se imprime el output de un filosofo que debiese ir antes que el de otro por ejemplo:

filosofo 3 "recibiendo palillo derecho"
filosofo 4 "enviando palillo izquierdo"

intuitivamente el filosofo 4 envio antes el palillo que el filosofo 3 lo recibiera y esto efectivamente ocurre asi
es solo en el print que esto se desordenada, cualquier duda mirar el codigo.

para correr el programa lo siguiente es necesario

docker, docker-compose

crear un archivo .env con lo siguiente
F1_IP=172.18.0.2
F2_IP=172.18.0.3
F3_IP=172.18.0.4
F4_IP=172.18.0.5
F5_IP=172.18.0.6
EAT_AMOUNT = 5

donde eatamount es el parametro de la cantidad de veces que cada filosofo debe comer.

las ips debiesen ser intercambiables ya que estas parametrizadas, pero no lo probé con otras.

existe un ejemplo del archivo .env en ejemploenv.txt copiar y pegar en .env en la raiz del repo

luego correr:

docker-compose up --build

listo.

El codigo es enredado, asi que lo explico a continuacion:
Cada contenedor alverga las subrutinas de mesero y filosofo

mesero y filosofo se comunican entre si por medio de canales eat done hungry

ademas existen otras 2 rutinas:

sender, reciever

donde sender se encarga de recibir mensajes dirigidos a otros meseros y enviarlos via udp
esta comunicacion se hace por medio de canales donde mesero envia a traves de un canal a sender los mensajes a ser enviados

reciever está escuchando constantemente y dirige los mensajes recibidos a traves de los canales correspondiente a mesero

cuando un filosofo come EAT_AMOUNT veces, manda un mensaje a sus vecinos que el ya termino, para terminar el programa un contenedor debe recibir
su señal de termino + la señal de termino de sus dos vecinos.

se siguio el algoritmo de las diapos.

cuando un mesero recibe la señal de que su filosofo tiene hambre el flujo es el siguiente:

tengo mis 2 palillos? entonces comer
no los tengo: pidelos y quedate esperando para luego comer
si alguien mientras tanto te pide tus palillos marca requested como True
luego de comer revisar requested l and r y envia a quien corresponda su palillo
(tambien se usar isDirty, etc.)

Cada filosofo espera una cantidad aleatoria entre 1 y 4 segundos de sleep asi que a menos que se ponga una seed cada ejecucion será distinta

tarea individual:

- ITAN FELZENSZTEIN 19639775
