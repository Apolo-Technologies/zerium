Pod::Spec.new do |spec|
  spec.name         = 'Gabt'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/apolo-technologies/zerium'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS Zerium Client'
  spec.source       = { :git => 'https://github.com/apolo-technologies/zerium.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Gzrm.framework'

	spec.prepare_command = <<-CMD
    curl https://gabtstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Gzrm.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
