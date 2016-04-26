'use strict';

module.exports = function (grunt) {
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		
		jshint: {
			options: {
				jshintrc: '.jshintrc'
			},
			all: [
				'Gruntfile.js',
				'dev/**/*.js', 
				'!dev/bower_components/**'
			]
		}, // jshint

		less: {
			dev: {
				files: { 
					'dev/assets/css/main.css': ['dev/assets/less/main.less'],
				},
				options: {
					compress: false
				}
			},

			prod: {
				files: {
					// Here it would be possible to create a custom bootstrap 
					// compile script that compiles only the bootstrap components
					// we actually need
					'prod/assets/css/app.min.css': [
						'dev/bower_components/bootstrap/less/bootstrap.less',
						'dev/bower_components/font-awesome/less/font-awesome.less',
						'dev/bower_components/toastr/toastr.less',
						'dev/bower_components/vex/css/vex.css',
						'dev/bower_components/vex/css/vex-theme-plain.css',
						'dev/bower_components/pickadate/lib/themes/classic.css',
						'dev/bower_components/pickadate/lib/themes/classic.date.css',
						'dev/bower_components/pickadate/lib/themes/classic.time.css',
						'dev/bower_components/bootstrap-select/less/bootstrap-select.less',
						'dev/assets/less/main.less'
					]
				},
				options: {
					compress: true,
					cleancss: true
				}
			}
		}, // less

		concat: {
			prod: {
				files: {
					'tmp/dependencies.js': [
						'dev/assets/js/ie10-viewport-bug-workaround.js',
						'dev/bower_components/jquery/dist/jquery.js',
						'dev/bower_components/bootstrap/js/transition.js',
						'dev/bower_components/bootstrap/js/collapse.js',
						'dev/bower_components/bootstrap/js/dropdown.js',
						'dev/bower_components/vex/js/vex.combined.min.js',
						'dev/bower_components/pickadate/lib/picker.js',
						'dev/bower_components/pickadate/lib/picker.date.js',
						'dev/bower_components/pickadate/lib/picker.time.js',
						'dev/bower_components/typeahead.js/dist/typeahead.bundle.js',
						'dev/bower_components/angular/angular.js',
						'dev/bower_components/angular-route/angular-route.js',
						'dev/bower_components/angular-timer/dist/angular-timer.min.js',
						'dev/bower_components/angular-cookies/angular-cookies.min.js',
						'dev/bower_components/angular-ui-bootstrap/src/modal/modal.js',
						'dev/bower_components/angular-ui-bootstrap/src/transition/transition.js',
						'dev/bower_components/toastr/toastr.js',
						'dev/bower_components/lodash/lodash.js',
						'dev/bower_components/momentjs/min/moment.min.js',
						'dev/bower_components/moment-duration-format/lib/moment-duration-format.js',
						'dev/bower_components/moment-timezone/builds/moment-timezone-with-data.min.js',
						'dev/ng-components/version/version.js',
						'dev/ng-components/version/version-directive.js',
						'dev/ng-components/version/interpolate-filter.js',
						'dev/ng-components/filters.js',
						'dev/bower_components/webfontloader/webfontloader.js',
						'dev/bower_components/bootstrap-select/js/bootstrap-select.js'
					],
					'tmp/app.js': [
						'tmp/dependencies.js',
						'dev/ng-main.js',
						'dev/ng-modules/activation/activation.js',
						'dev/ng-modules/activations/activations.js',
						'dev/ng-modules/api/api.js',
						'dev/ng-modules/bookings/bookings.js',
						'dev/ng-modules/coupons/coupons.js',
						'dev/ng-modules/coworking/coworking.js',
						'dev/ng-modules/coworking/coworking-product.js',
						'dev/ng-modules/coworking/coworking-purchase.js',
						'dev/ng-modules/dashboard/dashboard.js',
						'dev/ng-modules/invoices/invoices.js',
						'dev/ng-modules/login/login.js',
						'dev/ng-modules/mainmenu/mainmenu.js',
						'dev/ng-modules/machines/machines.js',
						'dev/ng-modules/machine/machine.js',
						'dev/ng-modules/memberships/memberships.js',
						'dev/ng-modules/membership/membership.js',
						'dev/ng-modules/priceunit/priceunit.js',
						'dev/ng-modules/productlist/productlist.js',
						'dev/ng-modules/randomtoken/randomtoken.js',
						'dev/ng-modules/reservation/reservation.js',
						'dev/ng-modules/reservations/reservations.js',
						'dev/ng-modules/reservations/toggle.js',
						'dev/ng-modules/space/space.js',
						'dev/ng-modules/spacepurchase/spacepurchase.js',
						'dev/ng-modules/spaces/spaces.js',
						'dev/ng-modules/tutoring/purchase.js',
						'dev/ng-modules/settings/settings.js',
						'dev/ng-modules/tutoring/tutor.js',
						'dev/ng-modules/tutoring/tutoring.js',
						'dev/ng-modules/user/user.js',
						'dev/ng-modules/users/users.js'
					]
				}
			}
		}, // concat

		uglify: {
			prod: {
				options: {
					compress: {
        		drop_console: true
      		}
				},
				files: {
					'prod/assets/js/modernizr.min.js': 'dev/bower_components/modernizr/modernizr.js',
					'prod/assets/js/app.min.js': 'tmp/app.js',
				}
			}
		}, // uglify

		processhtml: {
			prod: {
				options: {
					process: true
				}, 
				files: {
					'prod/index.html': ['dev/index.html']
				}
			}
		}, // processhtml

		copy: {
			prod: { // copy Angular template files to production dir
				files: [

				// Copy admin files
				{
					src: 'dev/bower_components/bootstrap/fonts/*',
					dest: 'prod/assets/fonts/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/bower_components/font-awesome/fonts/*',
					dest: 'prod/assets/fonts/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/coupons/coupons.html', 
					dest: 'prod/ng-modules/coupons/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/coworking/coworking.html', 
					dest: 'prod/ng-modules/coworking/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/coworking/coworking-product.html', 
					dest: 'prod/ng-modules/coworking/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/coworking/coworking-purchase.html', 
					dest: 'prod/ng-modules/coworking/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/dashboard/dashboard.html', 
					dest: 'prod/ng-modules/dashboard/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/dashboard/user-item.html', 
					dest: 'prod/ng-modules/dashboard/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/login/login.html', 
					dest: 'prod/ng-modules/login/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/mainmenu/mainmenu.html', 
					dest: 'prod/ng-modules/mainmenu/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/user/user.html', 
					dest: 'prod/ng-modules/user/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/activation/activation.html', 
					dest: 'prod/ng-modules/activation/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/activations/activations.html', 
					dest: 'prod/ng-modules/activations/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/machines/machines.html', 
					dest: 'prod/ng-modules/machines/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/machine/machine.html', 
					dest: 'prod/ng-modules/machine/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/memberships/memberships.html', 
					dest: 'prod/ng-modules/memberships/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/membership/membership.html', 
					dest: 'prod/ng-modules/membership/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/bookings/bookings.html', 
					dest: 'prod/ng-modules/bookings/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/invoices/invoices.html', 
					dest: 'prod/ng-modules/invoices/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/users/users.html', 
					dest: 'prod/ng-modules/users/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/priceunit/priceunit.html', 
					dest: 'prod/ng-modules/priceunit/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/productlist/productlist.html', 
					dest: 'prod/ng-modules/productlist/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/reservation/reservation.html', 
					dest: 'prod/ng-modules/reservation/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/reservations/reservations.html', 
					dest: 'prod/ng-modules/reservations/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/reservations/toggle.html', 
					dest: 'prod/ng-modules/reservations/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/settings/settings.html', 
					dest: 'prod/ng-modules/settings/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/space/space.html', 
					dest: 'prod/ng-modules/space/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/spacepurchase/spacepurchase.html', 
					dest: 'prod/ng-modules/spacepurchase/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/spaces/spaces.html', 
					dest: 'prod/ng-modules/spaces/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/tutoring/purchase.html', 
					dest: 'prod/ng-modules/tutoring/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/tutoring/tutor.html', 
					dest: 'prod/ng-modules/tutoring/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/ng-modules/tutoring/tutoring.html', 
					dest: 'prod/ng-modules/tutoring/',
					expand: true,
					flatten: true
				}, {
					src: 'dev/assets/img/*',
					dest: 'prod/assets/img/',
					filter: 'isFile',
					expand: true,
					flatten: true
				}]
			}
		}, // copy

		clean: {
			preprod: ['prod/**'],
			postprod: ['tmp/**']
		}, // clean

		watch: {
			js: {
				files: ['dev/**/*.js'],
				tasks: ['jshint'],
				options: {livereload: true, atBegin: true}
			},
			css: {
				files: ['dev/**/*.less'], 
				tasks: ['less:dev'],
				options: {livereload: true, atBegin: true}
			}
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
	grunt.registerTask('default', ['jshint', 'less:dev']);
	grunt.registerTask('dev', ['watch']);
	grunt.registerTask('prod', ['clean:preprod', 'jshint', 'less:prod', 'concat:prod', 'uglify:prod', 'processhtml:prod', 'copy:prod', 'clean:postprod']);
};
