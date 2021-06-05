FROM nginx:latest

# 설정 파일을 Docker 내부로 복사
COPY ./conf.d/default.conf /etc/nginx/conf.d/default.conf

# Docker Host로 포트 번호 노출
EXPOSE 80

# 백그라운드에서 실행
ENTRYPOINT ["nginx","-g","daemon off;"]