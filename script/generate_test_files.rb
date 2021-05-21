require 'faker'
require 'fileutils'

def print_usage()
  puts "usage: ruby generate_test_files.rb path/to/output"
end

def generate_files(path)
  paragraphs = generate_paragraphs 200
  write_content_files path, paragraphs, 350
end

def generate_paragraphs(num_paragraphs, num_sentences_min: 2, num_sentences_max: 10)
  paragraphs = []

  (0...num_paragraphs).each do
    num_sentences = Random.rand(num_sentences_max - num_sentences_min) + num_sentences_min
    paragraphs << Faker::Lorem.paragraph(sentence_count: num_sentences)
  end

  paragraphs
end

def write_content_files(path, paragraphs, num_files)
  content_path = File.join(path, 'content', 'items')
  FileUtils.mkdir_p(content_path) unless File.exist? content_path

  (1..num_files).each do |i|
    filepath = File.join(content_path, "test_content_#{i}.md")
    File.open(filepath, 'w') do |f|
      f.write("{{ title = '#{Faker::Book.title}' }}\n")
      f.write("\# {{: title}}\n"
      )
      (0...10).each do
        p = paragraphs[Random.rand(paragraphs.size)]
        bold = Random.rand(2) == 1

        if bold
          f.write("**#{p}**\n")
        else
          f.write("__#{p}__\n")
        end
      end
    end
  end
end

def write_pagination_template(path)
  template_path = File.join(path, 'templates')
  FileUtils.mkdir_p(template_path) unless File.exist? template_path

  filepath = File.join(template_path, "pagination_template.html")

  File.open(filepath, 'w') do |f|
    f.write("<p>The current page is {{: curPage }}</p>\n")
    f.write("{{for item in content}}\n")
    f.write("<div>{{: item._href }}</div>\n")
    f.write("{{end}}")
  end
end

def write_page_files(path)
  pages_path = File.join(path, 'pages')
  FileUtils.mkdir_p(pages_path) unless File.exist? pages_path

  write_content_pagination_page pages_path
end

def write_content_pagination_page(path)
  filepath = File.join(path, 'paged_content.html')
  File.open(filepath, 'w') do |f|
    f.write("{{: paginate(\"items\", \"pagination_template.html\", 10) }}\n")
    f.write("The cur page is: {{: curPage }}\n")
    f.write("{{for prevPage in pagesBefore(3)}}\n")
    f.write("<p>Before</p>\n")
    f.write("<a href=\"{{: prevPage._pageHref }}\">{{: prevPage._pageNum }}</a>\n")
    f.write("{{end}}\n")
    f.write("<strong>{{: curPage }}</strong>\n")
    f.write("{{for nextPage in pagesAfter(3)}}\n")
    f.write("<p>After</p>\n")
    f.write("<a href=\"{{: nextPage._pageHref }}\">{{: nextPage._pageNum }}</a>\n")
    f.write("{{end}}\n")
  end
end

if ARGV.length < 1
  print_usage
  return
end

output_path = ARGV[0]
generate_files output_path
write_pagination_template output_path
write_page_files output_path