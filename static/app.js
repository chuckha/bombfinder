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
		$('div.col').on('click', function (event) {
			var row = $(this).data('row');
			var col = $(this).data('col');
			ws.sendmsg('left', row, col);
			return false;
		});
		$('div.col').on('contextmenu', function(event) {
			var row = $(this).data('row');
			var col = $(this).data('col');
			ws.sendmsg('right', row, col);
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
    ws.sendmsg = function (type, row, col) {
      ws.send(JSON.stringify({'click': type, 'row': row, 'col': col}));
    };

    ws.onmessage = function (event) {
    	console.log(event);
    	ms.genBoard(JSON.parse(event['data']));
    };


}(window.ms = window.ms || {}, jQuery));