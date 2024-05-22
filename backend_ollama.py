from flask import Flask, render_template, request, jsonify, send_file
from flask_socketio import SocketIO, emit, disconnect
import os
import shutil
import subprocess
import ollama


app = Flask(__name__)
socketio = SocketIO(app, async_mode='threading')

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/yaml', methods=['POST'])
def generate():
    if 'file' not in request.files:
        return jsonify({'error': 'No file part'}), 400
    
    file = request.files['file']
    if file.filename == '':
        return jsonify({'error': 'No selected file'}), 400
    
    if file:
        # Ensure the micro directory exists
        micro_dir = 'micro'
        # Change directory to micro
        os.chdir(micro_dir)
        
        # Save the uploaded file as out.yaml
        file_path = 'out.yaml'
        file.save(file_path)
        
        # Run the command
        try:
            result = subprocess.run(['go', 'run', 'main.go', file_path], capture_output=True, text=True, check=True)
            output = result.stdout
        except subprocess.CalledProcessError as e:
            output = e.stderr
            return jsonify({'error': 'Command failed', 'output': output}), 500
        finally:
            # Change back to the original directory
            os.chdir('..')
        
        return jsonify({'message': 'File processed successfully', 'output': output}), 200



@app.route('/zip', methods=['GET'])
def download_zip():
    
    shutil.make_archive("out","zip","micro/output")
    try:
        return send_file('out.zip', as_attachment=True)
    except Exception as e:
        return str(e), 500
    finally:
        if os.path.exists('out.zip'):
            os.remove('out.zip')



@socketio.on('connect')
def handle_connect():
    print('Client connected')

@socketio.on('disconnect')
def handle_disconnect():
    print('Client disconnected')


@socketio.on('ask')
def handle_ask(data):
    user_input = data.get('question')
    yamlOutput = data.get('yamlOutput')
    print(user_input)
    print(yamlOutput)
    if not user_input:
        emit('error', {'message': 'No requirements provided'})
        return


    system_prompt = """
    You are given with a yaml file which will contain structure of a microservice. You will be asked to generate code in golang related to the yaml file. Just reply with code and do not explain the code. Do not print even ``` at begining and end of generated code.
    """
    query="YAML file:\n"+yamlOutput+"\n Question:\n"+user_input


    stream = ollama.generate(
        model="phi3",
        system=system_prompt,
        prompt=query,
        stream=True
    )

    for chunk in stream:
        emit('answer_chunk',{'answer':chunk['response']})

    disconnect()
    

@socketio.on('send_requirements')
def handle_requirements(data):
    user_input = data.get('requirements')
    if not user_input:
        emit('error', {'message': 'No requirements provided'})
        return

    print(f"Received requirements: {user_input}")


    system_prompt = """
You are a AI agent that converts microservice description in english to YAML file. Just output the yaml file and nothing else. Do not include code blocks using ``` too. The output should strictly follow the format given below and nothing else.
 Format of the yaml file must be as follows:

module: MenuService
port: 9000

database:
  provider: postgres
  url:  postgresql://postgres:l@localhost:3000/newdb2
  models:
    - table: Menu
      schema:
        Menuid: Int @id
        Name: String 
        Desc: String
        Availqty: Int

kafka: localhost:9092

endpoints: 
  - name: Addmenu
    path: /addmenu
    method: POST
    table: Menu
    kafka:
      topic: quantity
      type: consumer
    json:
        type: object
        properties:
            menuid: int 
            name: string
            desc: string
            availqty: int

  - name: Getmenu
    path: /menu
    method: GET
    table: Menu
    json:
        type: list
        properties:
            name: string
            desc: string

  - name: Getitem
    path: /menu/{id}
    method: GET
    table: Menu
    key:
      name: menuid
      type: int
    json:
        type: object
        properties:
            name: string
            desc: string
            availqty: int
    """


    stream = ollama.generate(
        model="phi3",
        system=system_prompt,
        prompt=str,
        stream=True
    )

    yaml_content = ""
    for chunk in stream:
        yaml_content += chunk['response']
        emit('yaml_chunk',{'yaml':chunk['response']})
    disconnect()


if __name__ == '__main__':
    socketio.run(app, debug=True)
