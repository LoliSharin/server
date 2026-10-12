#!/bin/bash

# Скрипт для переписи email-адресов во всех коммитах на sharinpetya000@gmail.com
# Репозитории, созданные с мая 2026

REPOS=(
  "server"
  "backend-api-core"
  "soft-report"
  "Outline-of-development-frontend"
  "Outline-of-development-backend"
  "Front-simulation"
  "business-_model"
)

GITHUB_USER="LoliSharin"
NEW_EMAIL="sharinpetya000@gmail.com"
NEW_NAME="Шарин Петр"

echo "🔧 Начинаем переписать email-адреса во всех коммитах..."
echo "📧 Целевой email: $NEW_EMAIL"
echo "👤 Имя автора: $NEW_NAME"
echo ""

for repo in "${REPOS[@]}"; do
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "📦 Обработка репозитория: $repo"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  # Создаём временную папку для каждого репо
  TEMP_DIR="/tmp/${repo}-fix"
  rm -rf "$TEMP_DIR"
  
  # Клонируем репозиторий
  echo "⬇️  Клонирование $repo..."
  git clone https://github.com/$GITHUB_USER/$repo.git "$TEMP_DIR"
  
  if [ $? -ne 0 ]; then
    echo "❌ Ошибка при клонировании $repo"
    continue
  fi
  
  cd "$TEMP_DIR"
  
  # Показываем текущие email-адреса
  echo "📋 Текущие email-адреса в коммитах:"
  git log --format='%ae' --all | sort | uniq
  echo ""
  
  # Переписываем email-адреса
  echo "🔄 Переписываем email-адреса..."
  git filter-branch -f --env-filter "
  if [ \"\$GIT_COMMITTER_EMAIL\" != \"$NEW_EMAIL\" ]; then
    export GIT_COMMITTER_EMAIL=\"$NEW_EMAIL\"
    export GIT_AUTHOR_EMAIL=\"$NEW_EMAIL\"
    export GIT_COMMITTER_NAME=\"$NEW_NAME\"
    export GIT_AUTHOR_NAME=\"$NEW_NAME\"
  fi
  " -- --all 2>&1 | grep -E "(Rewrite|refs|done|WARNING)"
  
  # Удаляем резервные копии
  echo "🧹 Очищаем резервные копии..."
  git for-each-ref --format='%(refname)' refs/original | xargs -r -n 1 git update-ref -d
  git reflog expire --expire=now --all
  git gc --prune=now --aggressive > /dev/null 2>&1
  
  # Определяем основную ветку
  MAIN_BRANCH=$(git rev-parse --abbrev-ref HEAD)
  
  # Протолкиваем изменения
  echo "📤 Отправляем изменения на GitHub..."
  git push --force origin $MAIN_BRANCH
  
  if [ $? -eq 0 ]; then
    echo "✅ $repo успешно обновлен!"
  else
    echo "❌ Ошибка при отправке $repo"
  fi
  
  echo ""
  
  # Возвращаемся в исходную директорию
  cd - > /dev/null
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Все репозитории обработаны!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Все коммиты теперь привязаны к email: $NEW_EMAIL"
echo "Они должны появиться в активности профиля в течение нескольких минут."
