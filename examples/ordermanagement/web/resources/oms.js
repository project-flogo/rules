
// MQTT client
var mqttClient;

function removeChildren(node) {
  if (node.hasChildNodes()) {
    while (node.childNodes.length >= 1) {
        node.removeChild(node.firstChild);
    }
  }
}

SendEventForm = function (id) {
  this.mId = id;
  this.mProperties = new Array();
  this.mDestinationPath;
  this.mEventPath;

  // add input field elements
  function _addField(parent, name, value, label, state) {
    var node;
    if ("hidden" == state) {
    } else {
        node = document.createElement('label');
        if (label) {
          node.appendChild(document.createTextNode(label + " "));
        }
        parent.appendChild(node);
        parent = node;
    }
    node = document.createElement('input');
    if ("hidden" == state) {
      node.setAttribute("type", "hidden");
    } else if ("disabled" == state) {
      node.setAttribute("disabled", "true");
    } else if ("readonly" == state) {
      node.setAttribute("readonly", "true");
    } else if ("password" == state) {
        node.setAttribute("type","password");
    }
    node.setAttribute("name", name);
    if (value) {
      node.setAttribute("value", value);
    }
    parent.appendChild(node);
  }

  // add multiline input element
  function _addMultiLineField(parent, name, rows, value, label) {
    var node;
    node = document.createElement('label');
    if (label) {
      node.appendChild(document.createTextNode(label + " "));
    }
    parent.appendChild(node);
    parent = node;
    node = document.createElement('textarea');
    if (rows) {
      node.setAttribute("rows", rows);
    }
    node.setAttribute("name", name);
    if (value) {
    node.appendChild(document.createTextNode(value));
    }
    parent.appendChild(node);
  }

  // add property field
  function _addPropertyField(parent, property, label, state) {
    var fieldName = property.name;
    if ("@" == fieldName.charAt(0)) {
      fieldName = "_" + fieldName.substr(1) + "_";
    }
    if (property.rows > 1) {
      if ("hidden" == state) {
        _addField(parent, fieldName, property.value, label, "hidden");
      } else if ("password" == state) {
        _addField(parent, fieldName, property.value, label, "password");
      } else {
        _addMultiLineField(parent, fieldName, property.rows, property.value, label, state);
      }
    } else {
      _addField(parent, fieldName, property.value, label, state);
    }
  }
  
  this.addProperty = function(name, value, rows, state) {
    var prop = new Object;
    prop.name = name;
    prop.value = value;
    prop.rows = rows;
    prop.state = state;
    this.mProperties[name] = prop;
  };

  this.targetName = "sendEventTarget";

  // builds the form element
  this.build = function() {
    var form = document.getElementById(this.mId);
    removeChildren(form);
    form.setAttribute("target", this.targetName);
    form.setAttribute("class", "sendEventForm " + form.className);
    form.setAttribute("action", this.mDestinationPath);
    form.setAttribute("onSubmit", "SendEventForm.onSubmit(this)");

    var fieldset = document.createElement('fieldset');
    form.appendChild(fieldset);
    
    var target = document.getElementById(this.targetName);
    if (target) {
    } else {
      target = document.createElement("iframe");
      target.setAttribute("name", this.targetName);
      target.setAttribute("id", this.targetName);
      target.setAttribute("class", "sendEventFormTarget");
      document.body.appendChild(target);
    }

    var e;
    e = document.createElement('legend');
    e.appendChild(document.createTextNode(this.mEventPath));
    fieldset.appendChild(e);

    var pre = document.createElement('pre');
    fieldset.appendChild(pre);

    var l = 0;
    for (var propName in this.mProperties){    
      if (propName.length > l) { l = propName.length; }
    }
    l;
    for (var propName in this.mProperties){    
      var prop = this.mProperties[propName];
      _addPropertyField(pre, prop, prop.name, prop.state);
      pre.appendChild(document.createTextNode("\n"));
    }

    var p = document.createElement('p');
    e = document.createElement('button');
    e.setAttribute("type", "submit");
    e.appendChild(document.createTextNode("Send \u2192 "));    
    var span = document.createElement('span');
    span.setAttribute("class", "sendEventFormDestination");
    span.appendChild(document.createTextNode(this.mDestinationPath));
    e.appendChild(span);
    p.appendChild(e);
    fieldset.appendChild(p);
  
    return fieldset;
  };

  this.getDestinationPath = function() {
    return this.mDestinationPath;
  };

  this.getEventPath = function() {
    return this.mEventPath;
  };

  this.setDestinationPath = function(path) {
    this.mDestinationPath = path;
  };

  this.setEventPath = function(path) {
    this.mEventPath = path;
  };

}

SendEventForm.get = function(id) {
    return new SendEventForm(id);
}

SendEventForm.getServer = function() {
  return SendEventForm.server;
}

// Process on action submit
SendEventForm.onSubmit = function(form) {
  event.stopPropagation();

  var propertiesJSON = "{";
  for (var i=0; i<form.elements.length; i++) {
    var element = form.elements[i]
    if (element.type == "text" && element.value != "") {
      propertiesJSON += ("\"" + element.name + "\":\"" + element.value + "\",")
    }
  }
  propertiesJSON = (propertiesJSON.substring(0, propertiesJSON.length-1) + "}");

  sendMsg(form.getAttribute("action"), propertiesJSON)
}

// Setup and connect to the server configured
SendEventForm.setServer = function(url) {
  SendEventForm.server = url;

  connectToMQTTBroker();
}

// Connect to MQTT broker
function connectToMQTTBroker() {
  var now = new Date();
  mqttClient = new Paho.MQTT.Client(SendEventForm.server, 8080, "clientId_"+now.getMilliseconds());
  mqttClient.onConnectionLost = onFailure;

  mqttClient.connect({onFailure:onFailure});
}

// called when the client connects
function sendMsg(topic, payload) {
  message = new Paho.MQTT.Message(payload);
  message.destinationName = topic;
  mqttClient.send(message);
}

// called when the client loses its connection
function onFailure(responseObject) {
  if (responseObject.errorCode !== 0) {
    console.log(responseObject.errorCode + " : " +responseObject.errorMessage);
  }
}


