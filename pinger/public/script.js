setInterval(function() {
		$.ajax({
			url: 'http://localhost:5000/get_list',
			type: 'GET',
			success: function(response){
				console.log(response);
			},
			error: function(error){
				console.log(error);
			}
		});
}, 10000);