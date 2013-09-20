;(function (ms, $) {
	var getClass = function(val) {
		if (!val || val == '-') {
			return '';
		}
		if (val === '⚑') {
			return 'flag';
		}
		return ['', 'one', 'two', 'three', 'four', 'five', 'six', 'seven', 'eight'][val];
	};

	var buildPlayer = function (playerObj) {
		var player = document.createElement('tr');
		var name = document.createElement('td');
		var active = document.createElement('td');
		var color = document.createElement('td');
		player.appendChild(name);
		player.appendChild(active);
		player.appendChild(color);
		name.innerHTML = playerObj['Name'];
		active.innerHTML = '☺';
		color.setAttribute('style', 'background-color: ' + playerObj['Color']);
		return player
	};

	var buildPlayers = function (listOfPlayers) {
		var i = 0,
			len = listOfPlayers.length,
			table = document.getElementById('playerTable');
		table.innerHTML = '<tr><th>Player</th><th>Active</th><th>Color</th></tr>';
		for (;i < len; i++) {
			table.appendChild(buildPlayer(listOfPlayers[i]));
		}
	};

	var buildCell = function (val, row, col) {
		var cell = document.createElement('div');
		cell.className = "col";
		if (val) {
			var d = document.createTextNode(val);
			cell.className += ' ' + getClass(val);
			cell.appendChild(d);
		}
		if (val != '' && val != '⚑' && val != '?') {
			cell.className += ' pressed';
		}
		if (val === '-') {
			cell.innerHTML = '';
		}
		cell.setAttribute('data-col', col);
		cell.setAttribute('data-row', row)
		return cell
	};

	var buildRow = function (listOfCells, rownum) {
		var i = 0,
			len = listOfCells.length,
			row = document.createElement('div');
		row.className = 'row';
		for (;i < len; i++) {
			row.appendChild(buildCell(listOfCells[i]['Val'], rownum, i));
		}
		row.firstChild.className += ' first';
		row.lastChild.className += ' last';
		return row;
	};

	var buildBoard = function (listOfRows) {
		var i = 0,
			len = listOfRows.length,
			board = document.createElement('div');
		board.id = 'board';
		for (;i < len; i++) {
			board.appendChild(buildRow(listOfRows[i], i));
		}
		board.firstChild.className += ' first';
		board.lastChild.className += ' last';
		return board;
	};

	ms.genBoard = function (obj) {
		var board = buildBoard(obj['Field']);
		var wrap = document.getElementById('game');
		wrap.innerHTML = '';
		wrap.appendChild(board);
		$('.col').on('click', function (event) {
			var row = $(this).data('row');
			var col = $(this).data('col');
			ws.move('left', row, col);
			return false;
		});
		$('.col').on('contextmenu', function(event) {
			var row = $(this).data('row');
			var col = $(this).data('col');
			ws.move('right', row, col);
			return false;
		});
	};
	ms.waitingMsg = function (obj) {
		// Looks for 'Required', 'Have' keys
		$('#startInfo').html('<p>Waiting for more players to join</p><p>Have ' + obj['Have'] + ' need ' + obj['Required'] + ' players.');
	}
	ms.updatePlayers = function (listOfPlayers) {
		// [
		//  {Color: "#aabbcc", Playing: true}
		// ]
		buildPlayers(listOfPlayers);
	}
	ms.invalidUsername = function (msg) {
		$('#errorusername').text(msg);
	};
	ms.setupUsername = function (username) {
		$('#usernamestuff').text('Playing as ' + username);
		$('#username-submit').off();
		$('#username').off();
	};
	/* events */
	$('#username-submit').on('click', function (event) {
		sendUsername();
	});
	$('#username').on('keypress', function (event) {
		if (event.keyCode === 13) {
			sendUsername();
		}
	});
	var sendUsername = function () {
		var username = $('#username').val();
		ws.send(username);
	}

    ws = new WebSocket("ws://localhost:8080/sweep");
    ws.onopen = function () {};
    ws.onclose = function () {
      // This is when the server closes the connection
    };
    ws.onerror = function () {};

    // Send things like {'type': 'leftclick', 'row': 3, 'col': 4}
    ws.move = function (type, row, col) {
    	var message = {
    		'Type': 'move',
    		'Value': {
    			'click': type, 
    			'row': row, 
    			'col': col
    		}
    	};
      ws.send(JSON.stringify(message));
    };

    ws.onmessage = function (event) {
    	console.log(event);

    	// just be ok with pings to see if we're still alive.
	   // if (event.data === 'ping') { return; }
    	var msg = JSON.parse(event['data']);
    	if (msg['Type'] === 'players') {
    		ms.updatePlayers(msg['Value']['Players']);
    	} else if (msg['Type'] === 'info') {
    		ms.waitingMsg(msg['Value']);
    	} else if (msg['Type'] === 'board') {
    		ms.genBoard(msg['Value']);
    	} else if (msg['Type'] === 'status') {
    		alert(msg['Value']);
    	} else if (msg['Type'] === 'InvalidUsername') {
    		ms.invalidUsername(msg['Value']);
    	} else if (msg['Type'] === 'ok') {
    		ms.setupUsername(msg['Value']);
    	}
    };
}(window.ms = window.ms || {}, jQuery));