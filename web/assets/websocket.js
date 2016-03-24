(function(){
  var socketUrl = "ws://"+document.location.host+"/Feed"
  var logLength = 20
  var hasConnected = false
  var pingInterval = 0
  var lastPing = new Date()
  var socketConnection

  var connectionStatus = function(s){
    var text = ""
    var highlight = ""

    if (s.readyState == 1){
      text = "Connected"
      highlight = "text-success"
      hasConnected = true
    } else if ([0,2,3].indexOf(s.readyState) >= 0){
      text = "Not Connected"
      highlight = "text-danger"
    }

    $('#connection-status').removeClass().addClass(highlight).text(text)
  }

  var connectSocket = function(){
    lastPing = new Date()
    var socket = new WebSocket(socketUrl)
    connectionStatus(socket)

    socket.onopen = function(){
      connectionStatus(socket)
    }

    socket.onclose = function(){
      connectionStatus(socket)
    }

    socket.onerror = function(e){
      connectionStatus(socket)
    }

    socket.onmessage = function(e){
      var msgBody = JSON.parse(e.data)
      msgBody.js_time = new Date(msgBody.time).toTimeString()
      msgBody.body_text = JSON.stringify(msgBody.body)

      $('#client-id').text(msgBody.client_id)
      $('#last-message').text(msgBody.js_time)
      $('#connected-clients').text(msgBody.connected_clients)

      $('#process-stats').html(Mustache.render(
        "Stats: Objects: {{heap_objects}}, Goroutines: {{goroutine_count}}, Memory: {{memory_usage}}",
        msgBody
      ))

      var newMessage = Mustache.render('<tr><td>{{message_type}}</td><td>{{body_text}}</td><td>{{js_time}}</td></tr>', msgBody)

      $('#message-log').prepend(newMessage)
      $('#message-log tr:nth-child(n+'+(logLength+1)+')').remove()

      var actions = {
        "ping": processPing,
        "helo": processHelo,
        "status_update": processStatus,
      }

      actions[msgBody.message_type](msgBody.body, msgBody)
    }

    function processHelo(b,d){
      pingInterval = d.ping_interval
      $('#ping-interval').text(pingInterval)

      if ( b != null ){
        updateStatus(b, false)
      }
    }

    function processPing(){
      lastPing = new Date()
    }

    function processStatus(b){
      updateStatus(b, true)
    }

    function updateStatus(status, autoplay){
      $('#current-status').text(status.text)

      $('#audio-player').children().remove()
      $('#audio-player').html(Mustache.render('<source src="{{audio_url}}" type="audio/mpeg" />', status))
      $('#audio-player').trigger('load')
      if ( autoplay ){
        $('#audio-player').trigger('play')
      }
    }

    return socket
  }

  window.onload = function(){
    socketConnection = connectSocket()
    setInterval(function(){
      if ( hasConnected == true && [0,1,2].indexOf(socketConnection.readyState) < 0 ){
        socketConnection = connectSocket()
      } else {
        var expectedPing = lastPing.getTime()+(pingInterval*1.25*1000)
        if ( pingInterval > 0 && new Date().getTime() > expectedPing ) {
          console.log("Missed ping")
          socketConnection.close()
        }
      }
    }, 2000)
  }
})()
