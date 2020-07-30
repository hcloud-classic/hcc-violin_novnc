# CentOS6 VNC 교체 (TigerVNC->RealVNC)

### 개요

CentOS 6 의 TigherVNC 가 violin-novnc 에서 VNC 접속시 인코딩 타입이 지원되지 않는 문제가 있어,

기존에 문제없이 작동하던 Debian 6 의 RealVNC 를 복사해와서 사용하는 방법을 제안해본다.



### 파일들 설명

- Xvnc4 : Debian 6 의 `/usr/bin/Xvnc4` 에 있던 파일로 RealVNC Server 파일이다.
- vncserver : Debian 6 의 `/usr/bin/vncserver`에 있던 파일로 Xvnc4 라이브러리를 통하여 VNC 서버 서비스를 구동하는 스크립트이다.
- X11.tar.gz : Debian 6 의 `/usr/share/fonts/X11/` 폴더를 압축해 놓은 파일로 Xvnc4 구동시 필요한 X11 폰트들이 담겨있다.
- hcc_init : VNC 서버 구동을 위한 hcc 스크립트 파일이다.



### 설치

1. X11.tar.gz 안의 `X11` 폴더를 `/usr/share/fonts/` 에 풀어준다.
2. `vncserver` 파일을 `/usr/bin/vncserver` 에 복사해준다.
3. `hcc_init` 파일을 `/etc/init.d/hcc_init` 에 복사해준다.
4. `/etc/init.d/hcc_init` 파일을 열어 `GEOMETRY` 값을 수정해준다.
   - ex) `GEOMETRY="800x600"`
5. `vncpasswd` 명령을 실행하여 VNC 접속 비밀번호를 설정해준다.
6. `hcc_init` 서비스를 시작해준다.
   - `service hcc_init start`