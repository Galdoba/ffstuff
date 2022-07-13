package actiontrailer

/*
# Трейлер HD:
MOVE \\nas\buffer\IN\FILE \\nas\buffer\IN\_IN_PROGRESS\ && \
fflite -r 25 -i \\nas\buffer\IN\_IN_PROGRESS\FILE \
-filter_complex "[0:a:0]aresample=48000,atempo=25/24[audio]" \
-map [audio]    @alac0 \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_AUDIORUS20.m4a \
-map 0:v:0      @crf10 \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_HD.mp4 \
&& MD \\nas\ROOT\IN\@TRAILERS\_DONE\NAME_TRL && MOVE \\nas\buffer\IN\_IN_PROGRESS\FILE \\nas\ROOT\IN\@TRAILERS\_DONE\NAME_TRL && MOVE \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_* \\nas\ROOT\EDIT\@trailers_temp\ && exit

# Трейлер HD DOWNSCALE:
MOVE \\nas\buffer\IN\FILE \\nas\buffer\IN\_IN_PROGRESS\ && \
fflite -r 25 -i \\nas\buffer\IN\_IN_PROGRESS\FILE \
-filter_complex "[0:a:0]aresample=48000,atempo=25/24[audio]; [0:v:0]scale=1920:-2,setsar=1/1,unsharp=3:3:0.3:3:3:0,pad=1920:1080:-1:-1[video_hd]" \
-map [audio]    @alac0 \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_AUDIORUS20.m4a \
-map [video_hd] @crf10 \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_HD.mp4 \
&& MD \\nas\ROOT\IN\@TRAILERS\_DONE\NAME_TRL && MOVE \\nas\buffer\IN\_IN_PROGRESS\FILE \\nas\ROOT\IN\@TRAILERS\_DONE\NAME_TRL && MOVE \\nas\buffer\IN\_IN_PROGRESS\NAME_TRL_* \\nas\ROOT\EDIT\@trailers_temp\ && exit

ПЛАН:
анализировать трейлер
выбрать интересующие стримы с которыми будем работать
формируем команду
	выбираем видео канал
		узнаем atempo
		узнаем размер
			узнаем надо ли даунскейлить
	выбираем звуковой канал
		выбираем язык
		выбираем канальность
прерываем если ошибка
перемещаем трейлер в папку inprogress
выполняем команду
перемещаем результаты в trailers_temp
перемещаем трейлер в trailers_done
завершаем программу
составить маршруты движения файлов
составить переменные для -filter_complex video
составить переменные для -filter_complex audio
запустить
*/
