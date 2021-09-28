# GOKHI
GO Kingdom Hearts Insider web scraper.  
  
This project aims to automate the process of downloading OST from khinsider.  
From thinking of an OST to listening to it is only 2 commands away !  
  
Here is a little demo of how to install and download an album:  
```bash
# Install gokhi
go get github.com/Ramoreik/gokhi

# Search for an album on khinsider
gokhi search -query '<insert-nostalgic-game-here>'

🔍 Query: \<nostalgic-game>
-- -- -- -- -- -- -- --
[0] ~+> \<nostalgic-game> The First
[1] ~+> \<nostalgic-game> The Second
[2] ~+> \<nostalgic-game> The Reboot
[3] ~+> \<nostalgic-game> The Remake
-- -- -- -- -- -- -- --
[?] use the 'download' subcommand to download an album.

# Download the chosen album 
gokhi download -album '<chosen-album-here' -download-path /home/somedude/Music/

🔍 Album to download: \<chosen-album>
🔍 Finding songs ...
🔍 Finding download links ...
🔍 Downloading found songs ...
🔽 Downloading /home/somedude/Music/\<chosen-album>/Song1.mp3
🔽 Downloading /home/somedude/Music/\<chosen-album>/Song2.mp3
🔽 Downloading /home/somedude/Music/\<chosen-album>/Song3.mp3
🔽 Downloading /home/somedude/Music/\<chosen-album>/Song4.mp3
🔽 Downloading /home/somedude/Music/\<chosen-album>/Song5.mp3
✅ Done, please check the specified download directory !

```
  
That's all, I don't think any other features need to be implemented really.  
I did this only to learn Golang and play with some goroutines.  
If you want to reuse or extend this projet, feel free to fork.  
  

