from ..models.conversation import Conversation
from ..chains.base import Handler
from ..enum.separation import Separation
class ContextInjectionHandler(Handler):
    def handle(self, convo: Conversation):
        convo.new_message = (
            "Remember previous discussion.\n" + convo.new_message + "\n After that can you generate the summary of your response after " + Separation.SEPRATION.value +" to separate the answer"
        )
        return super().handle(convo)


