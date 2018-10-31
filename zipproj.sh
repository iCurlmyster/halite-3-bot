rm halite.zip

cp src/main/MyBot.go MyBot.go

zip halite.zip -r MyBot.go src/helper/* src/hlt/* src/logic/* pkg/*

rm MyBot.go