Pod::Spec.new do |spec|
  spec.name         = 'Gzrm'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/abt/zerium'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS Zerium Client'
  spec.source       = { :git => 'https://github.com/abt/zerium.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Gzrm.framework'

	spec.prepare_command = <<-CMD
    curl https://gzrmstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Gzrm.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
