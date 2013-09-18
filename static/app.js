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

	var buildCell = function (val, row, col) {
		var cell = document.createElement('div');
		cell.className = "col";
		if (val) {
			var d = document.createTextNode(val);
			cell.className += ' ' + getClass(val);
			cell.appendChild(d);
		}
		if (val != '') {
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
		board.id = 'game';
		for (;i < len; i++) {
			board.appendChild(buildRow(listOfRows[i], i));
		}
		board.firstChild.className += ' first';
		board.lastChild.className += ' last';
		return board;
	};

	ms.genBoard = function (obj) {
		var board = buildBoard(obj['Field']);
		var wrap = document.getElementById('wrapper');
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
    ws.newGame = function () {
    	var message = {
    		'Type': 'newGame',
    		'Value': {}
    	};
    	ws.send(JSON.stringify(message));
    };

    ws.onmessage = function (event) {
    	console.log(event);
    	var msg = JSON.parse(event['data']);
    	if (msg['Type'] === 'board') {
    		ms.genBoard(msg['Value']);
    	}
    	if (msg['Type'] === 'status') {
    		alert(msg['Value']);
    	}
    };

    $('#newgame').on('click', function (event) {
    	ws.newGame();
    	return false;
    });
}(window.ms = window.ms || {}, jQuery));