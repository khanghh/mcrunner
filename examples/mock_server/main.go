package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	msg := []string{
		"A zombie was slain by Player123",
		"CreeperHunter mined diamond ore",
		"MineCrafter99 placed a block",
		"Server performance is normal",
		"Can't keep up! Is the server overloaded?",
		"Player456 just built an enormous castle with multiple towers, a moat, and even a drawbridge that actually functions in the game world, showcasing incredible creativity and dedication to the craft of Minecraft building.",
		"The Nether portal has been activated by AdventurerX, who ventured into the fiery dimension to collect blaze rods and fight off hordes of ghasts, skeletons, and piglins in a quest for rare enchanted gear and powerful potions.",
		"ExplorerY discovered a hidden village in the jungle biome, trading with villagers for emeralds and unlocking new recipes for advanced tools and armor.",
		"BuilderZ constructed a fully automated farm using pistons, observers, and redstone, harvesting crops efficiently without manual intervention.",
		"The End dimension was conquered by WarriorA, defeating the Ender Dragon and collecting dragon eggs for decorative purposes in the overworld.",
		"RedstoneMaster created a complex contraption with flying machines, TNT cannons, and automatic doors, demonstrating mastery over Minecraft's redstone mechanics.",
		"FarmerQ bred a herd of cows and chickens, establishing a sustainable food source and trading excess items with other players.",
	}
	count := 1
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	log.Println(msg[0])

	// Echo stdin to stdout
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			log.Println("received:", line)
		}
	}()

	for {
		select {
		case <-sigCh:
			log.Println("Shutting down server...")
			return
		case <-time.After(1 * time.Second):
			log.Println(msg[count%len(msg)])
			count++
		}
	}
}
