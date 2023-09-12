package main

import (
    "fmt"
    "os"

    "github.com/whisperapp/whisper/v2"
    "github.com/gotk3/gotk3/gtk"
)

const (
    // Spineta's Artaud CD color theme.
    artaudBackground = "#222222"
    artaudText = "#eeeeee"
    artaudButton = "#ffffff"
)

func main() {
    // Create a new Whisper client.
    client := whisper.NewClient()

    // Create a new recorder.
    recorder := client.NewRecorder()

    // Create a new window.
    window := gtk.NewWindow(gtk.WindowTypeToplevel)
    window.SetTitle("Whisper Recorder")
    window.Connect("destroy", func() {
        os.Exit(0)
    })

    // Set the window background color.
    window.SetBackgroundColor(artaudBackground)

    // Create a new box to contain the UI elements.
    box := gtk.NewBox(gtk.OrientationVertical, 10)

    // Create a button to start and stop recording.
    button := gtk.NewButtonWithLabel("Start Recording")
    button.Connect("clicked", func() {
        if recorder.IsRecording() {
            // Stop recording.
            recorder.Stop()
            button.SetLabel("Start Recording")
        } else {
            // Start recording.
            recorder.Start()
            button.SetLabel("Stop Recording")
        }
    })

    // Set the button color.
    button.SetBackgroundColor(artaudButton)

    // Create a text box to display the transcribed text.
    textBox := gtk.NewTextView()
    textBox.SetEditable(false)

    // Set the text box color.
    textBox.SetBackgroundColor(artaudBackground)
    textBox.SetTextColor(artaudText)

    // Add the UI elements to the box.
    box.Add(button)
    box.Add(textBox)

    // Add the box to the window.
    window.Add(box)

    // Show the window.
    window.ShowAll()

    // Start a new goroutine to listen for incoming audio data.
    go func() {
        for {
            // Read a chunk of audio data from the microphone.
            data, err := recorder.Read()
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            // Transcribe the audio data into text.
            text, err := client.Transcribe(data)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            // Display the transcribed text in the text box.
            textBox.SetText(text)
        }
    }()

    // Add a save button.
    saveButton := gtk.NewButtonWithLabel("Save")
    saveButton.Connect("clicked", func() {
        // Save the transcribed text to a file.
        fileName := "transcribed_text.txt"
        f, err := os.Create(fileName)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()

        text := textBox.GetText()
        f.WriteString(text)
    })
    box.Add(saveButton)

    // Start the main loop.
    gtk.Main()
}
