'use strict';

module.exports = function (grunt) {
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		
		less: {
			dev: {
				files: {
					'dev/global/global.css': ['dev/global/global.less'], 
					'dev/machines/main.css': ['dev/machines/main.less'],
					'dev/admin/main.css': ['dev/admin/main.less']
				},
				options: {
					compress: false
				}
			},
		}, // less

		watch: {
			css: {
				files: [
					'dev/**/*.less'
				],
				tasks: ['less:dev'],
				options: {
					livereload: true,
					atBegin: true
				}
			},
		} // watch
		
	});

	// Load tasks
	grunt.loadNpmTasks('grunt-contrib-clean');
	grunt.loadNpmTasks('grunt-contrib-jshint');
	grunt.loadNpmTasks('grunt-contrib-concat');
	grunt.loadNpmTasks('grunt-contrib-uglify');
	grunt.loadNpmTasks('grunt-notify');
	grunt.loadNpmTasks('grunt-contrib-watch');
	grunt.loadNpmTasks('grunt-text-replace');
	grunt.loadNpmTasks('grunt-processhtml');
	grunt.loadNpmTasks('grunt-contrib-less');
	grunt.loadNpmTasks('grunt-contrib-copy');

	// Register tasks
	grunt.registerTask('default', ['less']);
	grunt.registerTask('dev', ['watch']);  
};
