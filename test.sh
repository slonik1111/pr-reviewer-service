#!/bin/bash

echo
echo "=== 1. Добавить команду backend ==="
curl -s -X POST -H "Content-Type: application/json" \
  -d '{"team_name":"backend","members":[
        {"id":"u1","username":"Alice","is_active":true},
        {"id":"u2","username":"Bob","is_active":true}
      ]}' \
  http://localhost:8080/team/add
echo

echo
echo "=== 2. Получить команду backend ==="
curl -s http://localhost:8080/team/get?team_name=backend
echo


echo
echo "=== 3. Создать PR pr-1001 автор u1 ==="
curl -s -X POST -H "Content-Type: application/json" \
  -d '{"pull_request_id":"pr-1001","pull_request_name":"Add search","author_id":"u1"}' \
  http://localhost:8080/pullRequest/create
echo


echo
echo "=== 3.1. Показать все PR, где u1 является ревьювером ==="
curl -s http://localhost:8080/users/getReview?user_id=u1
echo


echo
echo "=== 3.2. Показать все PR, где u2 является ревьювером ==="
curl -s http://localhost:8080/users/getReview?user_id=u2
echo


echo
echo "=== 4. Пометить PR как MERGED ==="
curl -s -X POST -H "Content-Type: application/json" \
  -d '{"pull_request_id":"pr-1001"}' \
  http://localhost:8080/pullRequest/merge
echo


echo
echo "=== 5. Получить PR для ревьювера u2 ==="
curl -s http://localhost:8080/users/getReview?user_id=u2
echo


echo
echo "=== 6. Изменить активность пользователя u2 → false ==="
curl -s -X POST -H "Content-Type: application/json" \
  -d '{"user_id":"u2","is_active":false}' \
  http://localhost:8080/users/setIsActive
echo


echo
echo "=== 7. Переназначить ревьювера u2 → новый пользователь ==="
curl -s -X POST -H "Content-Type: application/json" \
  -d '{"pull_request_id":"pr-1001","old_user_id":"u2"}' \
  http://localhost:8080/pullRequest/reassign
echo


echo
echo "=== 8. Показать все PR, где u1 является ревьювером ==="
curl -s http://localhost:8080/users/getReview?user_id=u1
echo


echo
echo "=== 9. Показать все PR, где u2 является ревьювером ==="
curl -s http://localhost:8080/users/getReview?user_id=u2
echo

