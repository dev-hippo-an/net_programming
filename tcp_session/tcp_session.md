# TCP 데이터 스트림
## TCP 가 신뢰성 있는 이유
- tcp 는 패킷 손실과 패킷 수신 순서에 대한 문제를 잘 해결했다는 점에서 신뢰성이 있다.
- 데이터 전송 속도 조절을 통해 네트워크 상태가 변경되더라도 손실된 패킷을 최소로 유지하며 데이터를 빠르게 전송하도록 한다. (흐름 제어)
- 수신 패킷을 추적하고, 승인되지 않은 패킷은 필요에 따라 재전송한다.
- 흐름 제어, 패킷의 순차 처리, 재전송 등으로 tcp 는 패킷 손실 문제를 해결하고 데이터 전송에 대한 신뢰도를 높여 데이터에 집중할 수 있도록 한다.

## TCP 세션
- 3방향 핸드셰이크(두 노드간 패킷 전송)를 통해 tcp 연결을 한다. syn 패킷 -> syn 패킷 / ack 패킷 -> ack 패킷
- 연결 수립 후 두 노드는 데이터를 주고 받을 수 있는 상태가 된다.
- tcp 세션에서는 어느 한쪽에서 데이터를 전송하지 않으면 유휴 상태로 남으며 해당 세션은 메모리 낭비를 일으킨다.

## Go 언어에서 TCP 연결 처리
- net 패키지에는 tcp 관련 기능을 제공한다.
- 다만 연결을 적절하게 처리하는 부분은 코드상에서 처리해야 한다.

- net.Listen 함수 사용시 수신 연결 요청 처리가 가능한 tcp 서버를 작성할 수 있으며 이런 서버를 **리스너** 라고 한다.
- net.Dial 함수를 사용하여 리스너에 연결 시도를 할 수 있다.

### 하트비트
- 애플리케이션 계층에서 긴 유휴 시간을 가져야만 하는 경우 데드라인을 계속해서 뒤로 설정해 네트워크 연결을 지속시키기 위해 노드 간에 하트비트를 구현한다.
- 하트비트로 인해 네트워크상 장애를 빠르게 파악하고 연결을 재시도 할 수 있다.
