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
function generate() {
	let formData = {};
	
	formData.subject = document.getElementById('subject').value;
	formData.template = document.getElementById('template').value;

	formData.start_delim = document.getElementById('start_delim').value;
	formData.end_delim = document.getElementById('end_delim').value;

	formData.csv_delim = document.getElementById('csv_delim').value;

	let file = document.getElementById('data_file').value;
	console.log("debug: file passed:", file)

	formData.data = DataFile;
// TODO: Remove hardcode. Find the way to define filesize dy
	let generatedZip = new Uint8Array(2048); 
	generate_go(formData, generatedZip);
	console.log(result);
}

function populate_with_test_data() {
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
