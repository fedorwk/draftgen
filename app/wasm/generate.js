var DataFile;

function read_file(input) {
	let file = input.files[0];
	let reader = new FileReader();
	reader.readAsArrayBuffer(file);

  	reader.onload = function() {
		DataFile = new Uint8Array(reader.result) // necessary typecasting for go's syscall/js
  	};

  	reader.onerror = function() {
    	console.log(reader.error);
  	};
}

// TODO: wait until read_file() complete
async function generate() {
	let formData = {};
	
	formData.subject = document.getElementById('subject').value;
	formData.template = document.getElementById('template').value;
	formData.start_delim = document.getElementById('start_delim').value;
	formData.end_delim = document.getElementById('end_delim').value;
	formData.csv_delim = document.getElementById('csv_delim').value;

	formData.data = DataFile;

	// generate_go() is WASM function which generates result, stores it internally
	// returns resulting size in bytes, call retrieveResult_go to copy the result
	// to JS's Uint8Array
	let resultByteLength = generate_go(formData); 
	if (resultByteLength < 0) {
		console.log("error occured");
		return;
	}
	let dataBuffer = new Uint8Array(resultByteLength);
	// retrieveResult_go() is WASM function which retrieves generated file
	// from WASM and writes it to passed Uint8Array
	let bytesRetrieved = retrieveResult_go(dataBuffer);
	if (bytesRetrieved != resultByteLength) {
		console.log(`err: some bytes lost while copying: expected ${resultByteLength}, got ${bytesRetrieved}`);
		return;
	}
	sendFileToUser(dataBuffer);
}

function sendFileToUser(dataBuffer) {
	let blob = new Blob([dataBuffer]);
	let url = URL.createObjectURL(blob);
	let a = document.createElement('a');
	a.href = url;
	a.setAttribute("download","drafts.zip");
	document.body.append(a);
	a.click();
	a.remove();
	URL.revokeObjectURL(url);
}

function populateWithTestData() {
	document.getElementById('subject').value = 'TestSubject'
	document.getElementById('template').value = `Hello, {name}
You can buy {purchase} at {store} today`
	document.getElementById('start_delim').value = '{'
	document.getElementById('end_delim').value = '}'
	document.getElementById('csv_delim').value = ';'
	// TODO: implement
	// document.getElementById('data_file').setAttribute('value', '')
	alert("functionality incomplete, select testfile yourself")
}
