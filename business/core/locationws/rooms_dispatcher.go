package locationws

import "context"

var roomsDispatcher *RoomsDispatcher

type RoomsDispatcher struct {
	broadcastRoomMap map[string]*BroadcastRoom
}

func NewRoomsDispatcher() *RoomsDispatcher {
	roomsDispatcher = &RoomsDispatcher{
		broadcastRoomMap: make(map[string]*BroadcastRoom),
	}
	return roomsDispatcher
}

func AcquireBroadcastRoom(id string) *BroadcastRoom {
	if _, ok := roomsDispatcher.broadcastRoomMap[id]; !ok {
		roomsDispatcher.broadcastRoomMap[id] = newBroadcastRoom(id)
		ctx, cancel := context.WithCancel(context.Background())
		go roomsDispatcher.broadcastRoomMap[id].run(ctx)
		go roomsDispatcher.listenRoomClosed(roomsDispatcher.broadcastRoomMap[id], cancel)
	}
	return roomsDispatcher.broadcastRoomMap[id]
}

func (d *RoomsDispatcher) listenRoomClosed(room *BroadcastRoom, cancel context.CancelFunc) {
	<-room.channel.Closed
	cancel()
	delete(d.broadcastRoomMap, room.id)
}
