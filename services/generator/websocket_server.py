import logging

import websockets

logger = logging.getLogger(__name__)


class WebSocketServer:
    def __init__(self):
        self.clients = set()

    async def handler(self, websocket):
        self.clients.add(websocket)

        logger.info(
            "Client connected: %s",
            websocket.remote_address,
        )

        try:
            await websocket.wait_closed()
        finally:
            self.clients.discard(websocket)

            logger.info(
                "Client disconnected: %s",
                websocket.remote_address,
            )

    def broadcast(self, payload: str):
        if self.clients:
            websockets.broadcast(
                self.clients,
                payload,
            )

    @property
    def has_clients(self) -> bool:
        return len(self.clients) > 0
