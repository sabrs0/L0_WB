#!/bin/bash

# Флаг, указывающий на то, была ли уже произведена инициализация базы данных
INITIALIZED_FLAG="/var/lib/postgresql/data/.initialized"

# Проверяем, существует ли флаг инициализации
if [ ! -f "$INITIALIZED_FLAG" ]; then
  echo "Инициализация базы данных..."

  # Проверяем наличие файла init.sql внутри контейнера
  if [ -f "/docker-entrypoint-initdb.d/init.sql" ]; then
    echo " Выполняем загрузку данных из backup.sql"
    #psql -U postgres -d bees -f /docker-entrypoint-initdb.d/init.sql
  else
    echo "Файл init.sql не найден в контейнере."
    exit 1
  fi
  # Создаем флаг инициализации, чтобы при следующих запусках контейнера миграция не выполнялась
  touch "$INITIALIZED_FLAG"
else
  echo "База данных уже инициализирована, миграция не требуется."
fi
