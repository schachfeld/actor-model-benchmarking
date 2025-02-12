package coinbaseanalyzer


import spray.json._
import spray.json.DefaultJsonProtocol._

case class StartMessage(filename: String)
case class FileLineMessage(text: String)
case class WriteLineMessage(text: String)
case class LastMessage()
case class DoneMessage()
case class DoneMessageProduct(productId: String)

case class CBMessage(
                      channel: String,
                      client_id: String,
                      timestamp: String,
                      sequence_num: Int,
                      events: Seq[CBEvent]
                    )

case class CBEvent(
                    event_type: String,
                    product_id: String,
                    updates: Seq[CBUpdate]
                  )

case class CBUpdate(
                     side: String,
                     event_time: String,
                     price_level: String,
                     new_quantity: String
                   )

object CBJsonProtocol extends DefaultJsonProtocol {
  implicit val cbUpdateFormat: RootJsonFormat[CBUpdate] = jsonFormat4(CBUpdate)

  implicit val cbEventFormat: RootJsonFormat[CBEvent] = new RootJsonFormat[CBEvent] {
    override def write(obj: CBEvent): JsValue = JsObject(
      "type" -> JsString(obj.event_type), // Rename "event_type" to "type"
      "product_id" -> JsString(obj.product_id),
      "updates" -> obj.updates.toJson
    )

    override def read(json: JsValue): CBEvent = json.asJsObject.getFields("type", "product_id", "updates") match {
      case Seq(JsString(eventType), JsString(productId), JsArray(updates)) =>
        CBEvent(eventType, productId, updates.map(_.convertTo[CBUpdate]))
      case _ => throw DeserializationException("CBEvent expected")
    }
  }

  implicit val cbMessageFormat: RootJsonFormat[CBMessage] = jsonFormat5(CBMessage)
}

